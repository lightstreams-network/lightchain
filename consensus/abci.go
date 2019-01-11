package consensus

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/lightstreams-network/lightchain/database"
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
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

var bigZero = big.NewInt(0)

// maxTransactionSize is 32KB in order to prevent DOS attacks
const maxTransactionSize = 32768

type TendermintABCI struct {
	db           *database.Database
	checkTxState *state.StateDB
	rpcClient    *rpc.Client
	logger       tmtLog.Logger

	getCurrentDBState func() (*state.StateDB, error)
}

// Todo: assert this somehow
//var _ tmtAbciTypes.Application = TendermintABCI{}

func NewTendermintABCI(db *database.Database, client *rpc.Client, logger tmtLog.Logger) (*TendermintABCI, error) {
	txState, err := db.Ethereum().BlockChain().State()
	if err != nil {
		return nil, err
	}

	abci := &TendermintABCI{
		db:                db,
		rpcClient:         client,
		getCurrentDBState: db.Ethereum().BlockChain().State,
		checkTxState:      txState.Copy(),
		logger:			   logger,
	}
	
	return abci, nil
}

func (abci *TendermintABCI) InitEthState() error {
	return abci.db.InitEthState(abci.Receiver())
}

// Info returns information about the last height and app_hash to the tmtCfg engine
func (abci *TendermintABCI) Info(req tmtAbciTypes.RequestInfo) tmtAbciTypes.ResponseInfo {
	abci.logger.Info("TendermintABCI::Info()", "data", req)
	blockchain := abci.db.Ethereum().BlockChain()
	currentBlock := blockchain.CurrentBlock()
	height := currentBlock.Number()
	hash := currentBlock.Hash()

	abci.logger.Debug("Info", "height", height) // nolint: errcheck

	// This check determines whether it is the first time lightchain gets started.
	// If it is the first time, then we have to respond with an empty hash, since
	// that is what tmtCfg expects.
	if height.Cmp(bigZero) == 0 {
		return tmtAbciTypes.ResponseInfo{
			Data:             "ABCIEthereum",
			LastBlockHeight:  height.Int64(),
			LastBlockAppHash: []byte{},
		}
	}

	return tmtAbciTypes.ResponseInfo{
		Data:             "ABCIEthereum",
		LastBlockHeight:  height.Int64(),
		LastBlockAppHash: hash[:],
	}
}

// SetOption sets a configuration option
func (abci *TendermintABCI) SetOption(req tmtAbciTypes.RequestSetOption) tmtAbciTypes.ResponseSetOption {
	abci.logger.Info("TendermintABCI::SetOption()")
	return tmtAbciTypes.ResponseSetOption{Code: tmtAbciTypes.CodeTypeOK, Log: ""}
}

// InitChain initializes the validator set
func (abci *TendermintABCI) InitChain(req tmtAbciTypes.RequestInitChain) tmtAbciTypes.ResponseInitChain {
	abci.logger.Info("TendermintABCI::InitChain()")
	return tmtAbciTypes.ResponseInitChain{}
}

// CheckTx checks a transaction is valid but does not mutate the state
func (abci *TendermintABCI) CheckTx(txBytes []byte) tmtAbciTypes.ResponseCheckTx {
	abci.logger.Info("TendermintABCI::CheckTx()")
	tx, err := decodeRLP(txBytes)
	if err != nil {
		// nolint: errcheck
		abci.logger.Error("Received invalid transaction", "tx", tx.Hash().String())
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrEncodingError.Code), Log: err.Error()}
	}

	abci.logger.Info("Received valid transaction", "tx", tx.Hash().String()) // nolint: errcheck
	return abci.validateTx(tx)
}

// DeliverTx executes a transaction against the latest state
func (abci *TendermintABCI) DeliverTx(txBytes []byte) tmtAbciTypes.ResponseDeliverTx {
	abci.logger.Info("TendermintABCI::DeliverTx()")
	tx, err := decodeRLP(txBytes)
	if err != nil {
		// nolint: errcheck
		abci.logger.Info("DelivexTx: Received invalid transaction", "tx", tx, "err", err)
		return tmtAbciTypes.ResponseDeliverTx{Code: uint32(abciTypes.ErrEncodingError.Code), Log: err.Error()}
	}
	abci.logger.Info("DeliverTx: Received valid transaction", "tx", tx.Hash().String()) // nolint: errcheck

	res := abci.db.DeliverTx(tx)
	if res.IsErr() {
		// nolint: errcheck
		abci.logger.Error("DeliverTx: Error delivering tx to database db", "tx", tx.Hash().String(),
			"err", err)
		return res
	}

	abci.CollectTx(tx)
	return tmtAbciTypes.ResponseDeliverTx{Code: tmtAbciTypes.CodeTypeOK}
}

// BeginBlock starts a new Ethereum block
func (abci *TendermintABCI) BeginBlock(req tmtAbciTypes.RequestBeginBlock) tmtAbciTypes.ResponseBeginBlock {
	abci.logger.Info("TendermintABCI::BeginBlock()")
	abci.logger.Debug("BeginBlock") // nolint: errcheck

	// update the eth header with the tmtCfg header
	abci.db.UpdateHeaderWithTimeInfo(&req.Header)
	return tmtAbciTypes.ResponseBeginBlock{}
}

// EndBlock accumulates rewards for the validators and updates them
func (abci *TendermintABCI) EndBlock(req tmtAbciTypes.RequestEndBlock) tmtAbciTypes.ResponseEndBlock {
	abci.logger.Info("TendermintABCI::EndBlock()")
	//abci.db.AccumulateRewards(abci.strategy)
	return tmtAbciTypes.ResponseEndBlock{}
}

// Commits the block and returns a hash of the current state
func (abci *TendermintABCI) Commit() tmtAbciTypes.ResponseCommit {
	blockHash, err := abci.db.Commit(abci.Receiver())
	if err != nil {
		// nolint: errcheck
		abci.logger.Error("Error getting latest database state", "err", err)
	}

	ethState, err := abci.getCurrentDBState()
	if err != nil {
		abci.logger.Error("Error getting latest state", "err", err) // nolint: errcheck
	}

	abci.logger.Info("TendermintABCI::Commit()", "blockHash", blockHash.Hex())
	abci.checkTxState = ethState.Copy()
	// The app should respond to the Commit request with a byte array, which is the deterministic state root of the
	// application. It is included in the header of the next block. It can be used to provide easily verified
	// Merkle-proofs of the state of the application.
	return tmtAbciTypes.ResponseCommit{Data: blockHash.Bytes()}
}

// Query queries the state of the TendermintABCI
func (abci *TendermintABCI) Query(query tmtAbciTypes.RequestQuery) tmtAbciTypes.ResponseQuery {
	abci.logger.Info("TendermintABCI::Query()", "data", query)
	abci.logger.Debug("Query") // nolint: errcheck
	var in jsonRequest
	if err := json.Unmarshal(query.Data, &in); err != nil {
		return tmtAbciTypes.ResponseQuery{Code: uint32(abciTypes.ErrEncodingError.Code),
			Log: err.Error()}
	}
	var result interface{}
	if err := abci.rpcClient.Call(&result, in.Method, in.Params...); err != nil {
		return tmtAbciTypes.ResponseQuery{Code: uint32(abciTypes.ErrInternalError.Code),
			Log: err.Error()}
	}
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return tmtAbciTypes.ResponseQuery{Code: uint32(abciTypes.ErrInternalError.Code),
			Log: err.Error()}
	}
	return tmtAbciTypes.ResponseQuery{Code: tmtAbciTypes.CodeTypeOK, Value: resultBytes}
}

//-------------------------------------------------------

// validateTx checks the validity of a tx against the blockchain's current state.
// it duplicates the logic in database's tx_pool
func (abci *TendermintABCI) validateTx(tx *ethTypes.Transaction) tmtAbciTypes.ResponseCheckTx {
	abci.logger.Info("TendermintABCI::validateTx()", "data", tx.Hash().String())
	// Heuristic limit, reject transactions over 32KB to prevent DOS attacks
	if tx.Size() > maxTransactionSize {
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInternalError.Code),
			Log: core.ErrOversizedData.Error()}
	}

	var signer ethTypes.Signer = ethTypes.FrontierSigner{}
	if tx.Protected() {
		signer = ethTypes.NewEIP155Signer(tx.ChainId())
	}

	// Make sure the transaction is signed properly
	from, err := ethTypes.Sender(signer, tx)
	if err != nil {
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBaseInvalidSignature.Code),
			Log: core.ErrInvalidSender.Error()}
	}

	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if tx.Value().Sign() < 0 {
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBaseInvalidInput.Code),
			Log: core.ErrNegativeValue.Error()}
	}

	currentState := abci.checkTxState

	// Make sure the account exist - cant send from non-existing account.
	if !currentState.Exist(from) {
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBaseUnknownAddress.Code),
			Log: core.ErrInvalidSender.Error()}
	}

	// Check the transaction doesn't exceed the current block limit gas.
	gasLimit := abci.db.GasLimit()
	if gasLimit < 0 {
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInternalError.Code),
			Log: core.ErrGasLimitReached.Error()}
	}

	// Check if nonce is not strictly increasing
	nonce := currentState.GetNonce(from)
	if nonce != tx.Nonce() {
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBadNonce.Code),
			Log: fmt.Sprintf("Nonce not strictly increasing. Expected %d Got %d",
				nonce, tx.Nonce())}
	}

	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	currentBalance := currentState.GetBalance(from)
	if currentBalance.Cmp(tx.Cost()) < 0 {
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInsufficientFunds.Code),
			Log: fmt.Sprintf("Current balance: %s, tx cost: %s",
				currentBalance, tx.Cost())}
	}

	intrGas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil, true) // homestead == true
	if intrGas < 0 {
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBaseInsufficientFees.Code),
			Log: core.ErrIntrinsicGas.Error()}
	}
	if err != nil {
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInternalError.Code),
			Log: err.Error()}
	}

	// TODO: Evaluate usage of whitelist validation
	//isValid, err := abci.txHandler.IsValid(*tx)
	//if err != nil {
	//	return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInternalError.Code),
	//		Log: err.Error()}
	//}
	//if !isValid {
	//	msg := fmt.Sprintf("account %v not authorized to perform transaction %v", from.String(), tx.Hash().String())
	//	abci.logger.Info(msg)
	//	return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInternalError.Code),
	//		Log: err.Error()}
	//}

	// Update ether balances
	// amount + gasprice * gaslimit
	currentState.SubBalance(from, tx.Cost())
	// tx.To() returns a pointer to a common address. It returns nil
	// if it is a contract creation transaction.
	if to := tx.To(); to != nil {
		currentState.AddBalance(*to, tx.Value())
	}
	currentState.SetNonce(from, tx.Nonce()+1)

	return tmtAbciTypes.ResponseCheckTx{Code: tmtAbciTypes.CodeTypeOK}
}

//-------------------------------------------------------
// convenience methods for validators

// Receiver returns the receiving address based on the selected strategy
func (abci *TendermintABCI) Receiver() common.Address {
	//if abci.logger != nil {
	//	abci.logger.Info("TendermintABCI::Receiver()", "data", abci.strategy)
	//}

	return common.Address{}
}

// CollectTx invokes CollectTx on the strategy
func (abci *TendermintABCI) CollectTx(tx *ethTypes.Transaction) {
	abci.logger.Info("TendermintABCI::CollectTx()", "data", tx.Hash().String())
	abci.logger.Debug("CollectTx") // nolint: errcheck
}

// @TODO (ggarri): Refactor next ugly parsing and name and review
//func convertValidatorsToPointers(validators []tmtAbciTypes.Validator) []*tmtAbciTypes.Validator {
//	validatorPointers := []*tmtAbciTypes.Validator{}
//	for _, element := range validators {
//		validatorPointers = append(validatorPointers, &element)
//	}
//
//	return validatorPointers
//}

// RLP decode database transaction using go-database impl https://github.com/ethereum/go-ethereum/tree/v1.8.11/rlp
// TODO (ggarri): Align implementation with https://drive.google.com/file/d/11xB9ilEysXTar3samVE5Zki-QfPqejYj/view
func decodeRLP(txBytes []byte) (*ethTypes.Transaction, error) {
	tx := new(ethTypes.Transaction)
	rlpStream := rlp.NewStream(bytes.NewBuffer(txBytes), 0)
	if err := tx.DecodeRLP(rlpStream); err != nil {
		return nil, err
	}
	return tx, nil
}
