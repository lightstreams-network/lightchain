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

// format of query data
type jsonRequest struct {
	Method string          `json:"method"`
	ID     json.RawMessage `json:"id,omitempty"`
	Params []interface{}   `json:"params,omitempty"`
}

// maxTransactionSize is 32KB in order to prevent DOS attacks
const maxTransactionSize = 32768

type TendermintABCI struct {
	db           *database.Database
	checkTxState *state.StateDB
	ethRPCClient *rpc.Client
	logger       tmtLog.Logger

	getCurrentDBState func() (*state.StateDB, error)
	getCurrentBlock   func() *ethTypes.Block
}

var _ tmtAbciTypes.Application = &TendermintABCI{}

func NewTendermintABCI(db *database.Database, ethRPCClient *rpc.Client) (*TendermintABCI, error) {
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
	}

	return abci, nil
}

func (abci *TendermintABCI) InitEthState() error {
	return abci.db.InitEthState(abci.Receiver())
}

// Info returns information about the last height and app_hash to the tmtCfg engine
func (abci *TendermintABCI) Info(req tmtAbciTypes.RequestInfo) tmtAbciTypes.ResponseInfo {
	blockchain := abci.db.Ethereum().BlockChain()
	currentBlock := blockchain.CurrentBlock()
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

// SetOption sets a configuration option
func (abci *TendermintABCI) SetOption(req tmtAbciTypes.RequestSetOption) tmtAbciTypes.ResponseSetOption {
	abci.logger.Debug(fmt.Sprintf("Setting option key '%s' value '%s'", req.Key, req.Value))
	return tmtAbciTypes.ResponseSetOption{Code: tmtAbciTypes.CodeTypeOK, Log: ""}
}

// InitChain initializes the validator set
func (abci *TendermintABCI) InitChain(req tmtAbciTypes.RequestInitChain) tmtAbciTypes.ResponseInitChain {
	abci.logger.Debug(fmt.Sprintf("Initializing chain with ID '%s'", req.ChainId))
	return tmtAbciTypes.ResponseInitChain{}
}

// CheckTx checks a transaction is valid but does not mutate the state
func (abci *TendermintABCI) CheckTx(txBytes []byte) tmtAbciTypes.ResponseCheckTx {
	tx, err := decodeRLP(txBytes)
	if err != nil {
		abci.logger.Error("Received invalid transaction", "tx", tx.Hash().String())
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrEncodingError.Code), Log: err.Error()}
	}

	return abci.validateTx(tx)
}

// DeliverTx executes a transaction against the latest state
func (abci *TendermintABCI) DeliverTx(txBytes []byte) tmtAbciTypes.ResponseDeliverTx {
	tx, err := decodeRLP(txBytes)
	if err != nil {
		abci.logger.Info("Received invalid transaction", "tx", tx, "err", err)
		return tmtAbciTypes.ResponseDeliverTx{Code: uint32(abciTypes.ErrEncodingError.Code), Log: err.Error()}
	}

	res := abci.db.DeliverTx(tx)
	if res.IsErr() {
		abci.logger.Error("DeliverTx: Error delivering tx to database db", "tx", tx.Hash().String(), "err", err)
		return res
	}

	abci.CollectTx(tx)

	abci.logger.Info("TX delivered.", "tx", tx.Hash().String())

	return tmtAbciTypes.ResponseDeliverTx{Code: tmtAbciTypes.CodeTypeOK}
}

// BeginBlock starts a new Ethereum block
func (abci *TendermintABCI) BeginBlock(req tmtAbciTypes.RequestBeginBlock) tmtAbciTypes.ResponseBeginBlock {
	abci.logger.Debug(fmt.Sprintf("Beginning new block with hash '%s'", req.Hash))
	abci.db.UpdateHeaderWithTimeInfo(&req.Header)
	return tmtAbciTypes.ResponseBeginBlock{}
}

// EndBlock accumulates rewards for the validators and updates them
func (abci *TendermintABCI) EndBlock(req tmtAbciTypes.RequestEndBlock) tmtAbciTypes.ResponseEndBlock {
	abci.logger.Debug(fmt.Sprintf("Ending new block at height '%d'", req.Height))
	return tmtAbciTypes.ResponseEndBlock{}
}

// Commits the block and returns a hash of the current state
func (abci *TendermintABCI) Commit() tmtAbciTypes.ResponseCommit {
	rootHash := abci.getCurrentBlock().Root()
	blockHash, err := abci.db.Commit(abci.Receiver())
	if err != nil {
		abci.logger.Error("Error getting latest database state", "err", err)
		return tmtAbciTypes.ResponseCommit{Data: rootHash.Bytes()}
	}
	nextRootHash := abci.getCurrentBlock().Root()

	ethState, err := abci.getCurrentDBState()
	if err != nil {
		abci.logger.Error("Error getting latest state", "err", err)
		return tmtAbciTypes.ResponseCommit{Data: nextRootHash.Bytes()}
	}

	abci.logger.Info("Committing state", "blockHash", blockHash.Hex())

	abci.checkTxState = ethState.Copy()
	if err != nil {
		abci.logger.Error("Error committing latest state", "err", err)
		return tmtAbciTypes.ResponseCommit{Data: nextRootHash.Bytes()}
	}

	return tmtAbciTypes.ResponseCommit{Data: nextRootHash.Bytes()}
}

// Query queries the state of the TendermintABCI
func (abci *TendermintABCI) Query(query tmtAbciTypes.RequestQuery) tmtAbciTypes.ResponseQuery {
	abci.logger.Info("Querying state", "data", query)
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

// validateTx checks the validity of a tx against the blockchain's current state.
// it duplicates the logic in database's tx_pool
func (abci *TendermintABCI) validateTx(tx *ethTypes.Transaction) tmtAbciTypes.ResponseCheckTx {
	abci.logger.Info("Validating TX", "data", tx.Hash().String())

	if tx.Size() > maxTransactionSize {
	//if true {
		abci.logger.Error(core.ErrOversizedData.Error())
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInternalError.Code), Log: core.ErrOversizedData.Error()}
	}

	var signer ethTypes.Signer = ethTypes.FrontierSigner{}
	if tx.Protected() {
		signer = ethTypes.NewEIP155Signer(tx.ChainId())
	}

	// Make sure the transaction is signed properly
	from, err := ethTypes.Sender(signer, tx)
	if err != nil {
		abci.logger.Error(core.ErrInvalidSender.Error())
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBaseInvalidSignature.Code), Log: core.ErrInvalidSender.Error()}
	}

	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if tx.Value().Sign() < 0 {
		abci.logger.Error(core.ErrNegativeValue.Error())
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBaseInvalidInput.Code), Log: core.ErrNegativeValue.Error()}
	}

	currentState := abci.checkTxState

	if !currentState.Exist(from) {
		abci.logger.Error(core.ErrInvalidSender.Error())
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBaseUnknownAddress.Code), Log: core.ErrInvalidSender.Error()}
	}

	gasLimit := abci.db.GasLimit()
	if gasLimit < 0 {
		abci.logger.Error(core.ErrGasLimitReached.Error())
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInternalError.Code), Log: core.ErrGasLimitReached.Error()}
	}

	nonce := currentState.GetNonce(from)
	if nonce != tx.Nonce() {
		errMsg := fmt.Sprintf("Nonce not strictly increasing. Expected %d got %d", nonce, tx.Nonce())
		abci.logger.Error(errMsg)
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBadNonce.Code), Log: errMsg}
	}

	currentBalance := currentState.GetBalance(from)
	txCost := tx.Cost()
	if currentBalance.Cmp(txCost) < 0 {
		errMsg := fmt.Sprintf("Current balance: %s, tx cost: %s", currentBalance, tx.Cost())
		abci.logger.Error(errMsg)
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInsufficientFunds.Code), Log: errMsg}
	}

	intrGas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil, true) // homestead == true
	if intrGas < 0 {
		abci.logger.Error(core.ErrIntrinsicGas.Error())
		return tmtAbciTypes.ResponseCheckTx{
			Code: uint32(abciTypes.ErrBaseInsufficientFees.Code),
			Log: core.ErrIntrinsicGas.Error(),
		}
	}
	if err != nil {
		abci.logger.Error(err.Error())
		return tmtAbciTypes.ResponseCheckTx{
			Code: uint32(abciTypes.ErrInternalError.Code),
			Log: err.Error(),
		}
	}

	currentState.SubBalance(from, txCost)

	if to := tx.To(); to != nil {
		currentState.AddBalance(*to, tx.Value())
	}
	newAccountNonce := tx.Nonce() + 1
	currentState.SetNonce(from, newAccountNonce)

	abci.logger.Info("TX validated.", "tx", tx.Hash().String(), "cost", txCost.String(), "nonce", newAccountNonce)

	return tmtAbciTypes.ResponseCheckTx{Code: tmtAbciTypes.CodeTypeOK}
}

// Receiver returns the receiving address based on the selected strategy
func (abci *TendermintABCI) Receiver() common.Address {
	return common.Address{}
}

// CollectTx invokes CollectTx on the strategy
func (abci *TendermintABCI) CollectTx(tx *ethTypes.Transaction) {
	abci.logger.Info("Collecting TX", "data", tx.Hash().String())
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