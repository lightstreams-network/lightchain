package consensus

import (
	"encoding/json"
	"fmt"
	"math/big"
	"bytes"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/ethereum/go-ethereum/rpc"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	abciTypes "github.com/lightstreams-network/lightchain/consensus/types"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
	tmtLog "github.com/tendermint/tendermint/libs/log"
)

// maxTransactionSize is 32KB in order to prevent DOS attacks
const maxTransactionSize = 32768

// TendermintABCI is the main hook of application layer (blockchain) for connecting to consensus (Tendermint) using ABCI.
//
// Tendermint ABCI requires three connections to the application to handle the different message types.
// Connections:
//    Consensus Connection - InitChain, BeginBlock, DeliverTx, EndBlock, Commit
//    Mempool Connection - CheckTx
//    Info Connection - Info, SetOption, Query
//
// Flow:
//		1. BeginBlock
//		2. CheckTx
//	    3. DeliverTx
//		4. EndBlock
//		5. Commit
//		6. CheckTx (clean mempool from TXs not included in committed block)
//
// Tendermint runs CheckTx and DeliverTx concurrently with each other,
// though on distinct ABCI connections - the mempool connection and the consensus connection.
//
// Full ABCI specs: https://tendermint.com/docs/spec/abci/abci.html
type TendermintABCI struct {
	db           *database.Database
	checkTxState *state.StateDB
	ethRPCClient *rpc.Client
	logger       tmtLog.Logger
	metrics      *Metrics

	getCurrentDBState func() (*state.StateDB, error)
	getCurrentBlock   func() *ethTypes.Block
}

var _ tmtAbciTypes.Application = &TendermintABCI{}

func NewTendermintABCI(db *database.Database, ethRPCClient *rpc.Client, metrics *Metrics) (*TendermintABCI, error) {
	txState, err := db.Ethereum().BlockChain().State()
	if err != nil {
		return nil, err
	}

	abci := &TendermintABCI{
		db:                db,
		ethRPCClient:      ethRPCClient,
		getCurrentDBState: db.Ethereum().BlockChain().State,
		getCurrentBlock:   db.Ethereum().BlockChain().CurrentBlock,
		checkTxState:      txState.Copy(),
		logger:            log.NewLogger().With("engine", "consensus", "module", "ABCI"),
		metrics:           metrics,
	}

	return abci, nil
}

// InitChain is called upon genesis.
//
// Can be used to define validators set and consensus params on the application side.
//
// Response:
//     - If `ResponseInitChain.Validators` is empty, the initial validator set will be the `RequestInitChain.Validators`
func (abci *TendermintABCI) InitChain(req tmtAbciTypes.RequestInitChain) tmtAbciTypes.ResponseInitChain {
	abci.logger.Debug("Initializing chain", "chain_id", req.ChainId)

	return tmtAbciTypes.ResponseInitChain{}
}

// BeginBlock signals the beginning of a new block.
//
// The header contains the height, timestamp, and more - it exactly matches the Tendermint block header.
// The `req.LastCommitInfo` and `req.ByzantineValidators` attributes can be used to determine rewards and punishments
// for the validators.
//
// Response:
//		- Optional Key-Value tags for filtering and indexing
func (abci *TendermintABCI) BeginBlock(req tmtAbciTypes.RequestBeginBlock) tmtAbciTypes.ResponseBeginBlock {
	abci.logger.Debug("Beginning new block", "hash", req.Hash)
	abci.db.UpdateBlockState(&req.Header)

	return tmtAbciTypes.ResponseBeginBlock{}
}

// CheckTx validates a mempool transaction, prior to broadcasting or proposing.
//
// CheckTx should perform stateful but light-weight checks of the validity of
// the transaction (like checking signatures and account balances), but need
// not execute in full (like running a smart contract).
//
// The application should maintain a separate state to support CheckTx.
// This state can be reset to the latest committed state during Persist.
// Before calling Persist, Tendermint will lock and flush the mempool, ensuring
// that all existing CheckTx are responded to and no new ones can begin.
// After Persist, the mempool will rerun CheckTx for all remaining transactions,
// throwing out any that are no longer valid. Then the mempool will unlock
// and start sending CheckTx again.
//
// Response:
//		- Response code
//		- Optional Key-Value tags for filtering and indexing
func (abci *TendermintABCI) CheckTx(txBytes []byte) tmtAbciTypes.ResponseCheckTx {
	abci.metrics.CheckTxsTotal.Add(1)
	
	tx, err := decodeRLP(txBytes)
	if err != nil {
		abci.logger.Error("Received invalid transaction", "tx", tx.Hash().String())
		abci.metrics.CheckErrTxsTotal.Add(1, fmt.Sprint(abciTypes.ErrEncodingError.Code))
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrEncodingError.Code), Log: err.Error()}
	}

	abci.logger.Info("Checking TX", "hash", tx.Hash().String(), "nonce", tx.Nonce(), "cost", tx.Cost())

	if tx.Size() > maxTransactionSize {
		abci.logger.Error(core.ErrOversizedData.Error())
		abci.metrics.CheckErrTxsTotal.Add(1, fmt.Sprint(abciTypes.ErrInternalError.Code))
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInternalError.Code), Log: core.ErrOversizedData.Error()}
	}

	var signer ethTypes.Signer = ethTypes.FrontierSigner{}
	if tx.Protected() {
		signer = ethTypes.NewEIP155Signer(tx.ChainId())
	}

	from, err := ethTypes.Sender(signer, tx)
	if err != nil {
		abci.logger.Error(core.ErrInvalidSender.Error())
		abci.metrics.CheckErrTxsTotal.Add(1, fmt.Sprint(abciTypes.ErrBaseInvalidSignature.Code))
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBaseInvalidSignature.Code), Log: core.ErrInvalidSender.Error()}
	}

	if tx.Value().Sign() < 0 {
		abci.logger.Error(core.ErrNegativeValue.Error())
		abci.metrics.CheckErrTxsTotal.Add(1, fmt.Sprint(abciTypes.ErrBaseInvalidInput.Code))
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBaseInvalidInput.Code), Log: core.ErrNegativeValue.Error()}
	}

	if !abci.checkTxState.Exist(from) {
		abci.logger.Error(core.ErrInvalidSender.Error())
		abci.metrics.CheckErrTxsTotal.Add(1, fmt.Sprint(abciTypes.ErrBaseUnknownAddress.Code))
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBaseUnknownAddress.Code), Log: core.ErrInvalidSender.Error()}
	}

	if abci.checkTxState.GetNonce(from) != tx.Nonce() {
		errMsg := fmt.Sprintf("Nonce not strictly increasing. Expected %d got %d", abci.checkTxState.GetNonce(from), tx.Nonce())
		abci.logger.Error(errMsg)
		abci.metrics.CheckErrTxsTotal.Add(1, fmt.Sprint(abciTypes.ErrBadNonce.Code))
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBadNonce.Code), Log: errMsg}
	}

	intrinsicGas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil, true)
	if err != nil {
		abci.logger.Error(err.Error())
		abci.metrics.CheckErrTxsTotal.Add(1, fmt.Sprint(abciTypes.ErrInternalError.Code))
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInternalError.Code), Log: err.Error()}
	}

	if tx.Gas() < intrinsicGas {
		abci.logger.Error("TX gas is lower than intrinsic gas", "tx_gas", tx.Gas(), "intrinsic_gas", intrinsicGas)
		abci.metrics.CheckErrTxsTotal.Add(1, fmt.Sprint(abciTypes.ErrInsufficientGas.Code))
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInsufficientGas.Code), Log: err.Error()}
	}

	currentBalance := abci.checkTxState.GetBalance(from)
	txCost := tx.Cost()
	if currentBalance.Cmp(txCost) < 0 {
		errMsg := fmt.Sprintf("Current balance: %s, TX cost: %s", currentBalance, tx.Cost())
		abci.logger.Error(errMsg)
		abci.metrics.CheckErrTxsTotal.Add(1, fmt.Sprint(abciTypes.ErrInsufficientFunds.Code))
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInsufficientFunds.Code), Log: errMsg}
	}

	// Adjust From and To balances in order for the balance verification to be valid for next TX
	abci.checkTxState.SubBalance(from, txCost)
	if to := tx.To(); to != nil {
		abci.checkTxState.AddBalance(*to, tx.Value())
	}

	abci.checkTxState.SetNonce(from, tx.Nonce() + 1)

	abci.logger.Info("TX validated.", "hash", tx.Hash().String(), "state_nonce", abci.checkTxState.GetNonce(from))

	return tmtAbciTypes.ResponseCheckTx{Code: tmtAbciTypes.CodeTypeOK}
}

// DeliverTx executes the transaction against Ethereum block's work state.
//
// Response:
//		- If the transaction is valid, returns CodeType.OK
//		- Keys and values in Tags must be UTF-8 encoded strings
// 		  E.g: ("account.owner": "Bob", "balance": "100.0", "time": "2018-01-02T12:30:00Z")
func (abci *TendermintABCI) DeliverTx(txBytes []byte) tmtAbciTypes.ResponseDeliverTx {
	abci.metrics.DeliverTxsTotal.Add(1)
	tx, err := decodeRLP(txBytes)
	if err != nil {
		abci.metrics.DeliverErrTxsTotal.Add(1, fmt.Sprint(abciTypes.ErrEncodingError.Code))
		abci.logger.Info("Received invalid transaction", "hash", tx, "err", err)
		return tmtAbciTypes.ResponseDeliverTx{Code: uint32(abciTypes.ErrEncodingError.Code), Log: err.Error()}
	}

	abci.logger.Info("Delivering TX", "hash", tx.Hash().String(), "nonce", tx.Nonce(), "cost", tx.Cost(), "gas", tx.Gas(), "gas_price", tx.GasPrice())

	res := abci.db.ExecuteTx(tx)
	if res.IsErr() {
		abci.metrics.DeliverErrTxsTotal.Add(1, fmt.Sprint(res.Code))
		abci.logger.Error("Error delivering TX to DB", "hash", tx.Hash().String(), "err", res.Log)
		return res
	}

	abci.logger.Info("TX delivered.", "tx", tx.Hash().String())

	return tmtAbciTypes.ResponseDeliverTx{Code: tmtAbciTypes.CodeTypeOK}
}

// EndBlock signals the end of a block.
//
// An opportunity to propose changes to a validator set.
//
// Response:
// 		- Validator updates returned for block H
//			- apply to the NextValidatorsHash of block H+1
//			- apply to the ValidatorsHash (and thus the validator set) for block H+2
//			- apply to the RequestBeginBlock.LastCommitInfo (ie. the last validator set) for block H+3
//		- Consensus params returned for block H apply for block H+1
func (abci *TendermintABCI) EndBlock(req tmtAbciTypes.RequestEndBlock) tmtAbciTypes.ResponseEndBlock {
	abci.logger.Debug(fmt.Sprintf("Ending new block at height '%d'", req.Height))

	return tmtAbciTypes.ResponseEndBlock{}
}

// Commit persists the application state.
//
// Response:
//		- Return a Merkle root hash of the application state.
//	      It's critical that all application instances return the same hash. If not, they will not be able
// 		  to agree on the next block, because the hash is included in the next block!
func (abci *TendermintABCI) Commit() tmtAbciTypes.ResponseCommit {
	abci.metrics.CommitBlockTotal.Add(1)
	rootHash := abci.getCurrentBlock().Root()
	blockHash, err := abci.db.Persist(abci.RewardReceiver())
	if err != nil {
		abci.metrics.CommitErrBlockTotal.Add(1, "ErrGettingLastState")
		abci.logger.Error("Error getting latest database state", "err", err)
		return tmtAbciTypes.ResponseCommit{Data: rootHash.Bytes()}
	}
	nextRootHash := abci.getCurrentBlock().Root()

	ethState, err := abci.getCurrentDBState()
	if err != nil {
		abci.metrics.CommitErrBlockTotal.Add(1, "ErrGettingNextLastState")
		abci.logger.Error("Error getting next latest state", "err", err)
		return tmtAbciTypes.ResponseCommit{Data: nextRootHash.Bytes()}
	}

	abci.logger.Info("Committing state", "blockHash", blockHash.Hex())

	abci.checkTxState = ethState.Copy()

	return tmtAbciTypes.ResponseCommit{Data: nextRootHash.Bytes()}
}

// ResetBlockState resets the in-memory block's processing state.
func (abci *TendermintABCI) ResetBlockState() error {
	return abci.db.ResetBlockState(abci.RewardReceiver())
}

// RewardReceiver returns the receiving address based on the selected strategy
func (abci *TendermintABCI) RewardReceiver() common.Address {
	return common.Address{}
}

// Query for data from the application at current or past height.
func (abci *TendermintABCI) Query(query tmtAbciTypes.RequestQuery) tmtAbciTypes.ResponseQuery {
	abci.logger.Info("Querying state", "data", query)

	type jsonRequest struct {
		Method string          `json:"method"`
		ID     json.RawMessage `json:"id,omitempty"`
		Params []interface{}   `json:"params,omitempty"`
	}

	var in jsonRequest
	if err := json.Unmarshal(query.Data, &in); err != nil {
		return tmtAbciTypes.ResponseQuery{Code: uint32(abciTypes.ErrEncodingError.Code), Log: err.Error()}
	}

	var result interface{}
	if err := abci.ethRPCClient.Call(&result, in.Method, in.Params...); err != nil {
		return tmtAbciTypes.ResponseQuery{Code: uint32(abciTypes.ErrInternalError.Code), Log: err.Error()}
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return tmtAbciTypes.ResponseQuery{Code: uint32(abciTypes.ErrInternalError.Code), Log: err.Error()}
	}

	return tmtAbciTypes.ResponseQuery{Code: tmtAbciTypes.CodeTypeOK, Value: resultBytes}
}

// Info returns information about the last height and app_hash to the tmtCfg engine
func (abci *TendermintABCI) Info(req tmtAbciTypes.RequestInfo) tmtAbciTypes.ResponseInfo {
	currentBlock := abci.getCurrentBlock()
	height := currentBlock.Number()
	root := currentBlock.Root()

	abci.logger.Info("State info", "data", req, "height", height)

	// First boot-up
	if height.Cmp(big.NewInt(0)) == 0 {
		return tmtAbciTypes.ResponseInfo{
			Data:             "ABCIEthereum",
			LastBlockHeight:  height.Int64(),
			LastBlockAppHash: []byte{},
		}
	}

	return tmtAbciTypes.ResponseInfo{
		Data:             "ABCIEthereum",
		LastBlockHeight:  height.Int64(),
		LastBlockAppHash: root[:],
	}
}

// SetOption sets non-consensus critical application specific options.
//
// E.g. Key="min-fee", Value="100fermion" could set the minimum fee required
// for CheckTx (but not DeliverTx - that would be consensus critical).
func (abci *TendermintABCI) SetOption(req tmtAbciTypes.RequestSetOption) tmtAbciTypes.ResponseSetOption {
	abci.logger.Debug(fmt.Sprintf("Setting option key '%s' value '%s'", req.Key, req.Value))

	return tmtAbciTypes.ResponseSetOption{Code: tmtAbciTypes.CodeTypeOK, Log: ""}
}

// RLP decode database transaction using go-database impl https://github.com/ethereum/go-ethereum/tree/v1.8.11/rlp
// TODO: (ggarri): Align implementation with https://drive.google.com/file/d/11xB9ilEysXTar3samVE5Zki-QfPqejYj/view
func decodeRLP(txBytes []byte) (*ethTypes.Transaction, error) {
	tx := new(ethTypes.Transaction)
	rlpStream := rlp.NewStream(bytes.NewBuffer(txBytes), 0)
	if err := tx.DecodeRLP(rlpStream); err != nil {
		return nil, err
	}
	return tx, nil
}
