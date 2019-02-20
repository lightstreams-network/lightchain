package database

import (
	"github.com/ethereum/go-ethereum/core"
	"time"

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

	for obj := range db.ethTxsCh {
		db.logger.Debug("Captured NewTxsEvent from pool")
		db.metrics.BroadcastedTxsTotal.Add(1)
		for _, tx := range obj.Txs {
			if err := db.consAPI.BroadcastTx(*tx); err != nil {
				db.logger.Error("Error broadcasting tx", "err", err)
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
