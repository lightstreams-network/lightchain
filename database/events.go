package database

import (
	"time"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	txChanSize = 4096
)

// Transactions sent via the go-ethereum rpc need to be routed to tendermint.
//
// Listening to txs and forward to tendermint.
func (db *Database) txBroadcastLoop() {
	db.ethTxsCh = make(chan core.NewTxsEvent, txChanSize)
	db.ethTxSub = db.eth.TxPool().SubscribeNewTxsEvent(db.ethTxsCh)

	db.waitForTendermint()

	for {
		<-db.ethTxsCh
		db.logger.Debug("Captured NewTxsEvent from pool")

		queue, err := db.eth.TxPool().Pending()
		if err != nil {
			db.logger.Error("Error reading txPool pending queue", "error", err.Error())
			continue;
		}

		for from, txs := range queue {
			for _, tx := range txs {
				broadcastedTxNonce, ok := db.broadcastedTxCache[from]
				if !ok || broadcastedTxNonce < tx.Nonce() {
					db.logger.Debug("Broadcasting tx...", "from", from, "nonce", tx.Nonce())
					db.broadcastTx(from, *tx)
				} else {
					db.logger.Debug("Tx already broadcasted...", "from", from, "nonce", tx.Nonce())
				}
			}
		}
	}
}

// Wait for Tendermint RPC to open the socket and run http endpoint
func (db *Database) waitForTendermint() {
	for {
		status, err := db.consAPI.Status()
		if err == nil {
			break
		}

		db.logger.Info("Waiting for tendermint endpoint to start", "err", err, "result", status)
		time.Sleep(time.Second * 3)
	}

	db.logger.Info("Lightchain DB successfully connected to the Tendermint HTTP service.", "info", "Success")
}

func (db *Database) broadcastTx(from common.Address, tx types.Transaction) {
	if err := db.consAPI.BroadcastTx(tx); err != nil {
		db.metrics.BroadcastedErrTxsTotal.Add(1, err.Error())
		db.logger.Error("Error broadcasting tx", "err", err)
	} else {
		db.metrics.BroadcastedTxsTotal.Add(1)
		db.broadcastedTxCache[from] = tx.Nonce();
		db.logger.Debug("Broadcasted tx", "from", from, "nonce", tx.Nonce())
	}
}
