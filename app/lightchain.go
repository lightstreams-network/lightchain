package app

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/lightstreams-network/lightchain/ethereum"
	emtTypes "github.com/lightstreams-network/lightchain/types"

	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	abciTypesLegacy "github.com/lightstreams-network/lightchain/types/abci"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	tmLog "github.com/tendermint/tendermint/libs/log"
	"github.com/lightstreams-network/lightchain/core/transaction"
)

// format of query data
type jsonRequest struct {
	Method string          `json:"method"`
	ID     json.RawMessage `json:"id,omitempty"`
	Params []interface{}   `json:"params,omitempty"`
}

var bigZero = big.NewInt(0)

// maxTransactionSize is 32KB in order to prevent DOS attacks
const maxTransactionSize = 32768

type LightchainApplicationInterface interface {
	abciTypes.Application
	SetLogger(log tmLog.Logger)
}

// LightchainApplication implements an ABCI application
type LightchainApplication struct {
	// ethBackend handles the ethereum state machine
	// and wrangles other services started by an ethereum node (eg. tx pool)
	ethBackend *ethereum.Backend // ethBackend ethereum struct

	// a closure to return the latest current state from the ethereum blockchain
	getCurrentState func() (*state.StateDB, error)

	checkTxState *state.StateDB

	// an ethereum rpc client we can forward queries to
	rpcClient *rpc.Client

	// strategy for validator compensation
	strategy *emtTypes.Strategy

	logger tmLog.Logger

	// Handles all application logic related to transactions
	txHandler transaction.TxHandler
}

var _ LightchainApplicationInterface = (*LightchainApplication)(nil) // Verify that T implements I.

// @TODO (ggarri): Compare with "cosmos-sdk/baseapp/baseapp.go"
// Creates a fully initialised instance of LightchainApplication
func CreateLightchainApplication(ethBackend *ethereum.Backend,
	client *rpc.Client, strategy *emtTypes.Strategy, txHandler transaction.TxHandler) (*LightchainApplication, error) {

	txState, err := ethBackend.Ethereum().BlockChain().State()
	if err != nil {
		return nil, err
	}

	app := &LightchainApplication{
		ethBackend:      ethBackend,
		rpcClient:       client,
		getCurrentState: ethBackend.Ethereum().BlockChain().State,
		checkTxState:    txState.Copy(),
		strategy:        strategy,
		txHandler:       txHandler,
	}

	if err := app.ethBackend.InitEthState(app.Receiver()); err != nil {
		return nil, err
	}

	return app, nil
}

// SetLogger sets the logger for the lightchain application
func (app *LightchainApplication) SetLogger(log tmLog.Logger) {
	app.logger = log
}

// Info returns information about the last height and app_hash to the tendermint engine
func (app *LightchainApplication) Info(req abciTypes.RequestInfo) abciTypes.ResponseInfo {
	app.logger.Info("LightchainApplication::Info()", "data", req)
	blockchain := app.ethBackend.Ethereum().BlockChain()
	currentBlock := blockchain.CurrentBlock()
	height := currentBlock.Number()
	hash := currentBlock.Hash()

	app.logger.Debug("Info", "height", height) // nolint: errcheck

	// This check determines whether it is the first time lightchain gets started.
	// If it is the first time, then we have to respond with an empty hash, since
	// that is what tendermint expects.
	if height.Cmp(bigZero) == 0 {
		return abciTypes.ResponseInfo{
			Data:             "ABCIEthereum",
			LastBlockHeight:  height.Int64(),
			LastBlockAppHash: []byte{},
		}
	}

	return abciTypes.ResponseInfo{
		Data:             "ABCIEthereum",
		LastBlockHeight:  height.Int64(),
		LastBlockAppHash: hash[:],
	}
}

// SetOption sets a configuration option
func (app *LightchainApplication) SetOption(req abciTypes.RequestSetOption) abciTypes.ResponseSetOption {
	app.logger.Info("LightchainApplication::SetOption()", "data", req)
	app.logger.Debug("SetOption", "key", req.Key, "value", req.Value) // nolint: errcheck
	return abciTypes.ResponseSetOption{Code: abciTypes.CodeTypeOK, Log: ""}
}

// InitChain initializes the validator set
func (app *LightchainApplication) InitChain(req abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
	app.logger.Info("LightchainApplication::InitChain()", "data", req)
	app.logger.Debug("InitChain") // nolint: errcheck
	return abciTypes.ResponseInitChain{}
}

// CheckTx checks a transaction is valid but does not mutate the state
func (app *LightchainApplication) CheckTx(txBytes []byte) abciTypes.ResponseCheckTx {
	app.logger.Info("LightchainApplication::CheckTx()", "data", txBytes)
	tx, err := decodeRLP(txBytes)
	if err != nil {
		// nolint: errcheck
		app.logger.Info("CheckTx: Received invalid transaction", "tx", tx)
		return abciTypes.ResponseCheckTx{Code: uint32(abciTypesLegacy.ErrEncodingError.Code), Log: err.Error()}
	}
	app.logger.Info("CheckTx: Received valid transaction", "tx", tx) // nolint: errcheck

	return app.validateTx(tx)
}

// DeliverTx executes a transaction against the latest state
func (app *LightchainApplication) DeliverTx(txBytes []byte) abciTypes.ResponseDeliverTx {
	app.logger.Info("LightchainApplication::DeliverTx()", "data", txBytes)
	tx, err := decodeRLP(txBytes)
	if err != nil {
		// nolint: errcheck
		app.logger.Info("DelivexTx: Received invalid transaction", "tx", tx, "err", err)
		return abciTypes.ResponseDeliverTx{Code: uint32(abciTypesLegacy.ErrEncodingError.Code), Log: err.Error()}
	}
	app.logger.Info("DeliverTx: Received valid transaction", "tx", tx) // nolint: errcheck

	res := app.ethBackend.DeliverTx(tx)
	if res.IsErr() {
		// nolint: errcheck
		app.logger.Error("DeliverTx: Error delivering tx to ethereum ethBackend", "tx", tx,
			"err", err)
		return res
	}

	app.CollectTx(tx)
	return abciTypes.ResponseDeliverTx{Code: abciTypes.CodeTypeOK}
}

// BeginBlock starts a new Ethereum block
func (app *LightchainApplication) BeginBlock(req abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	app.logger.Info("LightchainApplication::BeginBlock()", "data", req)
	app.logger.Debug("BeginBlock") // nolint: errcheck

	// update the eth header with the tendermint header
	app.ethBackend.UpdateHeaderWithTimeInfo(&req.Header)
	return abciTypes.ResponseBeginBlock{}
}

// EndBlock accumulates rewards for the validators and updates them
func (app *LightchainApplication) EndBlock(req abciTypes.RequestEndBlock) abciTypes.ResponseEndBlock {
	app.logger.Info("LightchainApplication::EndBlock()", "data", req)
	app.logger.Debug("EndBlock", "height", req.Height) // nolint: errcheck
	//app.ethBackend.AccumulateRewards(app.strategy)
	return abciTypes.ResponseEndBlock{}
}

// Commits the block and returns a hash of the current state
func (app *LightchainApplication) Commit() abciTypes.ResponseCommit {
	app.logger.Debug("Commit") // nolint: errcheck
	blockHash, err := app.ethBackend.Commit(app.Receiver())
	if err != nil {
		// nolint: errcheck
		app.logger.Error("Error getting latest ethereum state", "err", err)
	}

	ethState, err := app.getCurrentState()
	if err != nil {
		app.logger.Error("Error getting latest state", "err", err) // nolint: errcheck
	}
	
	app.logger.Info("LightchainApplication::Commit()", "blockHash", blockHash.Bytes())
	app.checkTxState = ethState.Copy()
	return abciTypes.ResponseCommit{Data: blockHash.Bytes()}
}

// Query queries the state of the LightchainApplication
func (app *LightchainApplication) Query(query abciTypes.RequestQuery) abciTypes.ResponseQuery {
	app.logger.Info("LightchainApplication::Query()", "data", query)
	app.logger.Debug("Query") // nolint: errcheck
	var in jsonRequest
	if err := json.Unmarshal(query.Data, &in); err != nil {
		return abciTypes.ResponseQuery{Code: uint32(abciTypesLegacy.ErrEncodingError.Code),
			Log: err.Error()}
	}
	var result interface{}
	if err := app.rpcClient.Call(&result, in.Method, in.Params...); err != nil {
		return abciTypes.ResponseQuery{Code: uint32(abciTypesLegacy.ErrInternalError.Code),
			Log: err.Error()}
	}
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return abciTypes.ResponseQuery{Code: uint32(abciTypesLegacy.ErrInternalError.Code),
			Log: err.Error()}
	}
	return abciTypes.ResponseQuery{Code: abciTypes.CodeTypeOK, Value: resultBytes}
}

//-------------------------------------------------------

// validateTx checks the validity of a tx against the blockchain's current state.
// it duplicates the logic in ethereum's tx_pool
func (app *LightchainApplication) validateTx(tx *ethTypes.Transaction) abciTypes.ResponseCheckTx {
	app.logger.Debug("validateTx") // nolint: errcheck
	app.logger.Info("LightchainApplication::validateTx()", "data", tx)
	// Heuristic limit, reject transactions over 32KB to prevent DOS attacks
	if tx.Size() > maxTransactionSize {
		return abciTypes.ResponseCheckTx{Code: uint32(abciTypesLegacy.ErrInternalError.Code),
			Log: core.ErrOversizedData.Error()}
	}

	var signer ethTypes.Signer = ethTypes.FrontierSigner{}
	if tx.Protected() {
		signer = ethTypes.NewEIP155Signer(tx.ChainId())
	}

	// Make sure the transaction is signed properly
	from, err := ethTypes.Sender(signer, tx)
	if err != nil {
		return abciTypes.ResponseCheckTx{Code: uint32(abciTypesLegacy.ErrBaseInvalidSignature.Code),
			Log: core.ErrInvalidSender.Error()}
	}

	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if tx.Value().Sign() < 0 {
		return abciTypes.ResponseCheckTx{Code: uint32(abciTypesLegacy.ErrBaseInvalidInput.Code),
			Log: core.ErrNegativeValue.Error()}
	}

	currentState := app.checkTxState

	// Make sure the account exist - cant send from non-existing account.
	if !currentState.Exist(from) {
		return abciTypes.ResponseCheckTx{Code: uint32(abciTypesLegacy.ErrBaseUnknownAddress.Code),
			Log: core.ErrInvalidSender.Error()}
	}

	// Check the transaction doesn't exceed the current block limit gas.
	gasLimit := app.ethBackend.GasLimit()
	if gasLimit < 0 {
		return abciTypes.ResponseCheckTx{Code: uint32(abciTypesLegacy.ErrInternalError.Code),
			Log: core.ErrGasLimitReached.Error()}
	}

	// Check if nonce is not strictly increasing
	nonce := currentState.GetNonce(from)
	if nonce != tx.Nonce() {
		return abciTypes.ResponseCheckTx{Code: uint32(abciTypesLegacy.ErrBadNonce.Code),
			Log: fmt.Sprintf("Nonce not strictly increasing. Expected %d Got %d",
				nonce, tx.Nonce())}
	}

	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	currentBalance := currentState.GetBalance(from)
	if currentBalance.Cmp(tx.Cost()) < 0 {
		return abciTypes.ResponseCheckTx{Code: uint32(abciTypesLegacy.ErrInsufficientFunds.Code),
			Log: fmt.Sprintf("Current balance: %s, tx cost: %s",
				currentBalance, tx.Cost())}
	}

	intrGas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil, true) // homestead == true
	if intrGas < 0 {
		return abciTypes.ResponseCheckTx{Code: uint32(abciTypesLegacy.ErrBaseInsufficientFees.Code),
			Log: core.ErrIntrinsicGas.Error()}
	}
	if err != nil {
		return abciTypes.ResponseCheckTx{Code: uint32(abciTypesLegacy.ErrInternalError.Code),
			Log: err.Error()}
	}
	
	isValid, err := app.txHandler.IsValid(*tx)
	if err != nil {
		return abciTypes.ResponseCheckTx{Code: uint32(abciTypesLegacy.ErrInternalError.Code),
			Log: err.Error()}
	}
	if !isValid {
		msg := fmt.Sprintf("account %v not authorized to perform transaction %v", from.String(), tx.Hash().String())
		app.logger.Info(msg)
		return abciTypes.ResponseCheckTx{Code: uint32(abciTypesLegacy.ErrInternalError.Code),
			Log: err.Error()}
	}

	// Update ether balances
	// amount + gasprice * gaslimit
	currentState.SubBalance(from, tx.Cost())
	// tx.To() returns a pointer to a common address. It returns nil
	// if it is a contract creation transaction.
	if to := tx.To(); to != nil {
		currentState.AddBalance(*to, tx.Value())
	}
	currentState.SetNonce(from, tx.Nonce()+1)

	return abciTypes.ResponseCheckTx{Code: abciTypes.CodeTypeOK}
}

//-------------------------------------------------------
// convenience methods for validators

// Receiver returns the receiving address based on the selected strategy
func (app *LightchainApplication) Receiver() common.Address {
	if app.strategy != nil {
		app.logger.Debug("Receiver") // nolint: errcheck
		return app.strategy.Receiver()
	}

	if app.logger != nil {
		app.logger.Info("LightchainApplication::Receiver()", "data", app.strategy)
	}

	return common.Address{}
}

// CollectTx invokes CollectTx on the strategy
func (app *LightchainApplication) CollectTx(tx *ethTypes.Transaction) {
	app.logger.Info("LightchainApplication::CollectTx()", "data", tx)
	app.logger.Debug("CollectTx") // nolint: errcheck
	if app.strategy != nil {
		app.strategy.CollectTx(tx)
	}
}

// @TODO (ggarri): Refactor next ugly parsing and name and review
//func convertValidatorsToPointers(validators []abciTypes.Validator) []*abciTypes.Validator {
//	validatorPointers := []*abciTypes.Validator{}
//	for _, element := range validators {
//		validatorPointers = append(validatorPointers, &element)
//	}
//
//	return validatorPointers
//}

// RLP decode ethereum transaction using go-ethereum impl https://github.com/ethereum/go-ethereum/tree/v1.8.11/rlp
// TODO (ggarri): Align implementation with https://drive.google.com/file/d/11xB9ilEysXTar3samVE5Zki-QfPqejYj/view
func decodeRLP(txBytes []byte) (*ethTypes.Transaction, error) {
	tx := new(ethTypes.Transaction)
	rlpStream := rlp.NewStream(bytes.NewBuffer(txBytes), 0)
	if err := tx.DecodeRLP(rlpStream); err != nil {
		return nil, err
	}
	return tx, nil
}
