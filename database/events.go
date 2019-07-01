package database

import (
	"time"
	"github.com/ethereum/go-ethereum/core"
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

	go db.processTxQueueLoop()

	for obj := range db.ethTxsCh {
		db.logger.Debug("Captured NewTxsEvent from pool")
		for _, tx := range obj.Txs {
			db.logger.Debug("Adding to tx queue...", "nonce", tx.Nonce())
			db.txQueue.addTx(tx)
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

func (db *Database) processTxQueueLoop() {
	for {
		if !db.txQueue.isReady() {
			time.Sleep(time.Second)
			continue;
		}
		
		for tx := db.txQueue.popTx(); tx != nil; tx = db.txQueue.popTx() {
			db.logger.Debug("Broadcasting tx...", "nonce", tx.Nonce())
			if err := db.consAPI.BroadcastTx(*tx); err != nil {
				db.metrics.BroadcastedErrTxsTotal.Add(1, err.Error())
				db.logger.Error("Error broadcasting tx", "err", err)
			} else {
				db.logger.Debug("Broadcasted tx", "nonce", tx.Nonce())
				db.metrics.BroadcastedTxsTotal.Add(1)
			}
		}
	}
}
