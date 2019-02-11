package consensus

import (
	"github.com/ethereum/go-ethereum/core"
	"fmt"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
	abciTypes "github.com/lightstreams-network/lightchain/consensus/types"
)

// maxTransactionSize is 32KB in order to prevent DOS attacks
const maxTransactionSize = 32768

// CheckTx validates a mempool transaction, prior to broadcasting or proposing.
//
// Flow:
//		1. BeginBlock
//		2. CheckTx <- ******
//	    3. ExecuteTx
//		4. EndBlock
//		5. Commit
//		6. CheckTx <- ****** (clean mempool from TXs not included in committed block)
//
// CheckTx should perform stateful but light-weight checks of the validity of
// the transaction (like checking signatures and account balances), but need
// not execute in full (like running a smart contract).
//
// Tendermint runs CheckTx and ExecuteTx concurrently with each other,
// though on distinct ABCI connections - the mempool connection and the consensus connection, respectively.
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
//
func (abci *TendermintABCI) CheckTx(txBytes []byte) tmtAbciTypes.ResponseCheckTx {
	tx, err := decodeRLP(txBytes)
	if err != nil {
		abci.logger.Error("Received invalid transaction", "tx", tx.Hash().String())
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrEncodingError.Code), Log: err.Error()}
	}

	abci.logger.Info("Checking TX", "hash", tx.Hash().String(), "nonce", tx.Nonce(), "cost", tx.Cost())

	if tx.Size() > maxTransactionSize {
		abci.logger.Error(core.ErrOversizedData.Error())
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInternalError.Code), Log: core.ErrOversizedData.Error()}
	}

	var signer ethTypes.Signer = ethTypes.FrontierSigner{}
	if tx.Protected() {
		signer = ethTypes.NewEIP155Signer(tx.ChainId())
	}

	from, err := ethTypes.Sender(signer, tx)
	if err != nil {
		abci.logger.Error(core.ErrInvalidSender.Error())
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBaseInvalidSignature.Code), Log: core.ErrInvalidSender.Error()}
	}

	if tx.Value().Sign() < 0 {
		abci.logger.Error(core.ErrNegativeValue.Error())
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBaseInvalidInput.Code), Log: core.ErrNegativeValue.Error()}
	}

	if !abci.checkTxState.Exist(from) {
		abci.logger.Error(core.ErrInvalidSender.Error())
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBaseUnknownAddress.Code), Log: core.ErrInvalidSender.Error()}
	}

	if abci.checkTxState.GetNonce(from) != tx.Nonce() {
		errMsg := fmt.Sprintf("Nonce not strictly increasing. Expected %d got %d", abci.checkTxState.GetNonce(from), tx.Nonce())
		abci.logger.Error(errMsg)
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrBadNonce.Code), Log: errMsg}
	}

	intrinsicGas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil, true)
	if err != nil {
		abci.logger.Error(err.Error())
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInternalError.Code), Log: err.Error()}
	}

	if tx.Gas() < intrinsicGas {
		abci.logger.Error("TX gas is lower than intrinsic gas", "tx_gas", tx.Gas(), "intrinsic_gas", intrinsicGas)
		return tmtAbciTypes.ResponseCheckTx{Code: uint32(abciTypes.ErrInsufficientGas.Code), Log: err.Error()}
	}

	currentBalance := abci.checkTxState.GetBalance(from)
	txCost := tx.Cost()
	if currentBalance.Cmp(txCost) < 0 {
		errMsg := fmt.Sprintf("Current balance: %s, TX cost: %s", currentBalance, tx.Cost())
		abci.logger.Error(errMsg)
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