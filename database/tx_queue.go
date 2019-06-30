package database

import (
	"time"
	"sync"
	
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"sort"
)

const awaitingTxsTimeInMs = 1 * time.Second

type txQueue struct {
	lastUpdate time.Time
	txs        []*ethTypes.Transaction
	mtx        sync.Mutex
}

func (q *txQueue) addTx(tx *ethTypes.Transaction) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	q.lastUpdate = time.Now()
	q.txs = append(q.txs, tx)
	sort.Slice(q.txs, func(i, j int) bool { 
		return q.txs[i].Nonce() > q.txs[j].Nonce() 
	})
}

func (q *txQueue) popTx() (*ethTypes.Transaction) {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if len(q.txs) == 0 {
		return nil
	}

	var tx *ethTypes.Transaction
	tx, q.txs = q.txs[len(q.txs)-1], q.txs[:len(q.txs)-1]
	return tx
}

func (q *txQueue) isReady() bool {
	if len(q.txs) == 0 {
		return false
	}

	isReady := q.lastUpdate.Add(awaitingTxsTimeInMs).Unix() < time.Now().Unix()
	return isReady
}