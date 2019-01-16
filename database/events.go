package database

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
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
		log.Info("Captured NewTxsEvent from pool")

		for _, tx := range obj.Txs {
			log.Info("New TX", "data", tx.Hash().String())
			if err := db.consensusAPI.BroadcastTx(*tx); err != nil {
				log.Error("Error broadcasting tx", "err", err)
			}
		}
	}
}

// Wait for Tendermint RPC to open the socket and run http endpoint
func (db *Database) waitForTendermint() {
	for {
		status, err := db.consensusAPI.Status()
		if err == nil {
			break
		}

		log.Info("Waiting for tendermint endpoint to start", "err", err, "result", status)
		time.Sleep(time.Second * 3)
	}

	log.Info("Lightchain DB successfully connected to the Tendermint HTTP service.", "info", "Success")
}