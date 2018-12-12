package ethereum

import (
	"bytes"
	"github.com/ethereum/go-ethereum/core"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	rpcTypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcClient "github.com/tendermint/tendermint/rpc/lib/client"
	"time"
)

const (
	txChanSize = 4096
)

//----------------------------------------------------------------------
// Transactions sent via the go-ethereum rpc need to be routed to tendermint

// Listening to txs and forward to tendermint
func (b *Backend) txBroadcastLoop() {
	b.txsCh = make(chan core.NewTxsEvent, txChanSize)
	b.txSub = b.ethereum.TxPool().SubscribeNewTxsEvent(b.txsCh)

	waitForTendermint(b.client)

	log.Info("Block for tendermint endpoint to start", "info", "Success")

	for obj := range b.txsCh {
		log.Info("Captured NewTxsEvent from pool")
		for _, tx := range obj.Txs {
			log.Info("New Trx", "data", tx.Hash().String())
			if err := b.BroadcastTx(tx); err != nil {
				log.Error("Broadcast error", "err", err)
			}
		}
	}
}

// BroadcastTx broadcasts a transaction to tendermint core
func (b *Backend) BroadcastTx(tx *ethTypes.Transaction) error {
	result := new(rpcTypes.SyncInfo)
	log.Info("Broadcasts a transaction to Tendermint core")
	buf := new(bytes.Buffer)
	if err := tx.EncodeRLP(buf); err != nil {
		log.Error("Broadcasts a transaction to Tendermint core", "err", err)
		return err
	}
	// TODO: Review
	params := map[string]interface{}{
		"tx": buf.Bytes(),
	}
	
	// Broadcast_tx_sync will return with the result of running the transaction through CheckTx
	// alternatively we could use `broadcast_tx_async` which will return right away without waiting to hear if the 
	// transaction is even valid 
	_, err := b.client.Call("broadcast_tx_sync", params, result)
	return err
}

//----------------------------------------------------------------------
// Wait for Tendermint RPC to open the socket and run http endpoint
// Cosmos-sdk correlative call in cosmos-sdk/tests/util.go
func waitForTendermint(client rpcClient.HTTPClient) {
	result := new(rpcTypes.ResultStatus)
	for {
		_, err := client.Call("status", map[string]interface{}{}, &result)
		if err == nil {
			break
		}
		log.Info("Waiting for tendermint endpoint to start", "err", err, "result", result)
		time.Sleep(time.Second * 3)
	}
}
