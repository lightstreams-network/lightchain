package consensus

import (
	"encoding/json"
	"fmt"
	"math/big"
	"bytes"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
	tmtLog "github.com/tendermint/tendermint/libs/log"

	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus/metrics"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/tendermint/tendermint/crypto"
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
// 		1. BeginBlock
// 		2. CheckTx
// 	    3. DeliverTx
// 		4. EndBlock
// 		5. Commit
// 		6. CheckTx (clean mempool from TXs not included in committed block)
//
// Tendermint runs CheckTx and DeliverTx concurrently with each other,
// though on distinct ABCI connections - the mempool connection and the consensus connection.
//
// Full ABCI specs: https://tendermint.com/docs/spec/abci/abci.html
type TendermintABCI struct {
	db           *database.Database
	checkTxState *state.StateDB
	ethRPCClient *rpc.Client

	logger         tmtLog.Logger
	metrics        metrics.Metrics
	curBlockHeader *tmtAbciTypes.Header

	getCurrentDBState   func() (*state.StateDB, error)
	getCurrentBlock     func() *ethTypes.Block
	getValidatorAddress func(pubKey string) (common.Address, error)
}

var _ tmtAbciTypes.Application = &TendermintABCI{}

func NewTendermintABCI(db *database.Database, ethRPCClient *rpc.Client, metrics metrics.Metrics) (*TendermintABCI, error) {
	txState, err := db.Ethereum().BlockChain().State()
	if err != nil {
		return nil, err
	}

	abci := &TendermintABCI{
		db:                  db,
		ethRPCClient:        ethRPCClient,
		getCurrentDBState:   db.Ethereum().BlockChain().State,
		getCurrentBlock:     db.Ethereum().BlockChain().CurrentBlock,
		getValidatorAddress: db.Validators().ValidatorAddress,
		checkTxState:        txState.Copy(),
		logger:              log.NewLogger().With("engine", "consensus", "module", "ABCI"),
		metrics:             metrics,
	}

	err = abci.ResetBlockState()
	if err != nil {
		return nil, err
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
// 		- Optional Key-Value tags for filtering and indexing
func (abci *TendermintABCI) BeginBlock(req tmtAbciTypes.RequestBeginBlock) tmtAbciTypes.ResponseBeginBlock {
	abci.logger.Debug("Beginning new block", "hash", req.Hash)
	abci.db.UpdateBlockState(&req.Header)
	abci.curBlockHeader = &req.Header

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
// 		- Response code
// 		- Optional Key-Value tags for filtering and indexing
func (abci *TendermintABCI) CheckTx(txBytes []byte) tmtAbciTypes.ResponseCheckTx {
	abci.metrics.CheckTxsTotal.Add(1)

	tx, err := decodeRLP(txBytes)
	if err != nil {
		abci.logger.Error("Unable to decode RLP TX", "err", err.Error())
		abci.metrics.CheckErrTxsTotal.Add(1, "INVALID_TX")
		return tmtAbciTypes.ResponseCheckTx{Code: 1, Log: "INVALID_TX"}
	}

	abci.logger.Info("Checking TX", "hash", tx.Hash().String(), "nonce", tx.Nonce(), "cost", tx.Cost())

	var signer ethTypes.Signer = ethTypes.FrontierSigner{}
	if tx.Protected() {
		signer = ethTypes.NewEIP155Signer(tx.ChainId())
	}

	from, err := ethTypes.Sender(signer, tx)
	if err != nil {
		abci.logger.Error("Unable to retrieve TX sender", "err", err.Error())
		abci.metrics.CheckErrTxsTotal.Add(1, "INVALID_SENDER")
		return tmtAbciTypes.ResponseCheckTx{Code: 1, Log: "INVALID_SENDER"}
	}

	err = abci.doMempoolValidation(tx, from)
	if err != nil {
		abci.logger.Error(err.Error())
		abci.metrics.CheckErrTxsTotal.Add(1, err.Error())
		return tmtAbciTypes.ResponseCheckTx{Code: 1, Log: err.Error()}
	}

	// Additional custom validation aside of the Ethereum default mempool validation
	if tx.GasPrice().Uint64() < abci.db.Config().TxPool.PriceLimit {
		abci.logger.Error("TX gas price too low", "tx_gas_price", tx.GasPrice(), "required_tx_gas_price", abci.db.Config().TxPool.PriceLimit)
		abci.metrics.CheckErrTxsTotal.Add(1, "INSUFFICIENT_GAS_PRICE")
		return tmtAbciTypes.ResponseCheckTx{Code: 1, Log: "INSUFFICIENT_GAS_PRICE"}
	}

	currentBalance := abci.checkTxState.GetBalance(from)
	txCost := tx.Cost()
	if currentBalance.Cmp(txCost) < 0 {
		abci.logger.Error("Low balance", "balance", currentBalance.String(), "tx_cost", tx.Cost().String())
		abci.metrics.CheckErrTxsTotal.Add(1, "INSUFFICIENT_FUNDS")
		return tmtAbciTypes.ResponseCheckTx{Code: 1, Log: "INSUFFICIENT_FUNDS"}
	}

	// Adjust From and To balances in order for the balance verification to be valid for next TX
	abci.checkTxState.SubBalance(from, txCost)
	if to := tx.To(); to != nil {
		abci.checkTxState.AddBalance(*to, tx.Value())
	}

	abci.checkTxState.SetNonce(from, tx.Nonce()+1)

	abci.logger.Info("TX validated", "hash", tx.Hash().String(), "state_nonce", abci.checkTxState.GetNonce(from))

	return tmtAbciTypes.ResponseCheckTx{Code: tmtAbciTypes.CodeTypeOK}
}

func (abci *TendermintABCI) doMempoolValidation(tx *ethTypes.Transaction, from common.Address) (err error) {
	if tx.Size() > maxTransactionSize {
		return fmt.Errorf("MAX_TRANSACTION_SIZE")
	}

	if tx.Value().Sign() < 0 {
		return fmt.Errorf("INVALID_SIGNATURE")
	}

	if !abci.checkTxState.Exist(from) {
		return fmt.Errorf("UNKNOWN_ADDRESS")
	}

	if abci.checkTxState.GetNonce(from) != tx.Nonce() {
		abci.logger.Error(fmt.Sprintf("Nonce not strictly increasing. Expected %d got %d", abci.checkTxState.GetNonce(from), tx.Nonce()))
		return fmt.Errorf("BAD_NONCE")
	}

	intrinsicGas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil, true)
	if err != nil {
		return fmt.Errorf("INTRINSIC_GAS_UNKNOWN")
	}

	if tx.Gas() < intrinsicGas {
		return fmt.Errorf("INSUFFICIENT_INTRINSIC_GAS")
	}

	return nil
}

// DeliverTx executes the transaction against Ethereum block's work state.
//
// Response:
// 		- If the transaction is valid, returns CodeType.OK
// 		- Keys and values in Tags must be UTF-8 encoded strings
// 		  E.g: ("account.owner": "Bob", "balance": "100.0", "time": "2018-01-02T12:30:00Z")
func (abci *TendermintABCI) DeliverTx(txBytes []byte) tmtAbciTypes.ResponseDeliverTx {
	abci.metrics.DeliverTxsTotal.Add(1)
	tx, err := decodeRLP(txBytes)
	if err != nil {
		abci.logger.Error(err.Error())
		abci.metrics.DeliverErrTxsTotal.Add(1, "INVALID_TX")
		return tmtAbciTypes.ResponseDeliverTx{Code: 1, Log: "INVALID_TX"}
	}

	abci.logger.Info("Delivering TX", "hash", tx.Hash().String(), "nonce", tx.Nonce(), "cost", tx.Cost(), "gas", tx.Gas(), "gas_price", tx.GasPrice())

	res := abci.db.ExecuteTx(tx)
	if res.IsErr() {
		abci.logger.Error("Error delivering TX to DB", "hash", tx.Hash().String(), "err", res.Log)
		abci.metrics.DeliverErrTxsTotal.Add(1, "UNABLE_TO_DELIVER")
		return res
	}

	abci.logger.Info("TX delivered", "tx", tx.Hash().String())

	return tmtAbciTypes.ResponseDeliverTx{Code: tmtAbciTypes.CodeTypeOK}
}

// EndBlock signals the end of a block.
//
// An opportunity to propose changes to a validator set.
//
// Response:
// 		- Validator updates returned for block H
// 			- apply to the NextValidatorsHash of block H+1
// 			- apply to the ValidatorsHash (and thus the validator set) for block H+2
// 			- apply to the RequestBeginBlock.LastCommitInfo (ie. the last validator set) for block H+3
// 		- Consensus params returned for block H apply for block H+1
func (abci *TendermintABCI) EndBlock(req tmtAbciTypes.RequestEndBlock) tmtAbciTypes.ResponseEndBlock {
	abci.logger.Debug(fmt.Sprintf("Ending new block at height '%d'", req.Height))

	return tmtAbciTypes.ResponseEndBlock{}
}

// Commit persists the application state.
//
// Response:
// 		- Return a Merkle root hash of the application state.
// 	      It's critical that all application instances return the same hash. If not, they will not be able
// 		  to agree on the next block, because the hash is included in the next block!
func (abci *TendermintABCI) Commit() tmtAbciTypes.ResponseCommit {
	abci.metrics.CommitBlockTotal.Add(1)
	
	block, err := abci.db.Persist(abci.RewardReceiver())
	if err != nil {
		abci.logger.Error("Error getting latest database state", "err", err)
		abci.metrics.CommitErrBlockTotal.Add(1, "UNABLE_TO_PERSIST")
		panic(err)
	}

	ethState, err := abci.getCurrentDBState()
	if err != nil {
		abci.logger.Error("Error getting next latest state", "err", err)
		abci.metrics.CommitErrBlockTotal.Add(1, "ErrGettingNextLastState")
	}
	abci.checkTxState = ethState.Copy()

	abci.logger.Info("Block committed", "block", block.Hash().Hex(), "root", block.Root().Hex())

	return tmtAbciTypes.ResponseCommit{Data: block.Root().Bytes()}
}

// ResetBlockState resets the in-memory block's processing state.
func (abci *TendermintABCI) ResetBlockState() error {
	return abci.db.ResetBlockState(abci.RewardReceiver())
}

// RewardReceiver returns the receiving address based on the selected strategy
func (abci *TendermintABCI) RewardReceiver() common.Address {
	if abci.curBlockHeader == nil {
		abci.logger.Error("Missing block headers");
		return common.HexToAddress("") 
	}
	
	pubKeyAddr := crypto.Address(abci.curBlockHeader.GetProposerAddress())
	address, err := abci.getValidatorAddress(pubKeyAddr.String())
	if err != nil {
		abci.logger.Error("Cannot fetch validator rewarded address", "pubkey", pubKeyAddr.String(), "err", err.Error())
		return common.HexToAddress("")
	}
	
	abci.logger.Info("Block proposer validator is rewarded", "pubkey", pubKeyAddr.String(), "address", address.String())

	return address
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
		abci.logger.Error("unable to unmarshal query data", "err", err.Error())
		return tmtAbciTypes.ResponseQuery{Code: 1, Log: "INVALID_QUERY_REQ"}
	}

	var result interface{}
	if err := abci.ethRPCClient.Call(&result, in.Method, in.Params...); err != nil {
		abci.logger.Error("unable to call ETH RPC client", "err", err.Error())
		return tmtAbciTypes.ResponseQuery{Code: 1, Log: "ETH_RPC_CLIENT_UNAVAILABLE"}
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		abci.logger.Error("unable to unmarshal info response", "err", err.Error())
		return tmtAbciTypes.ResponseQuery{Code: 1, Log: "INVALID_QUERY_RES"}
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
