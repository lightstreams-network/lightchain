package txclient

import (
	"time"
)

const TxReceiptInterval = 1 * time.Second
const TxReceiptTimeout = 30 * time.Second

type TxConfig struct {
	gasLimit uint64
	txReceiptTimeout time.Duration
	txReceiptInterval time.Duration
}

func NewTxConfig(gasLimit uint64, receiptTimeout time.Duration, receiptInterval time.Duration) TxConfig {
	return TxConfig{
		gasLimit,
		receiptTimeout,
		receiptInterval,
	}
}

func NewTransferTxConfig() TxConfig {
	return TxConfig{
		21000,
		TxReceiptTimeout,
		TxReceiptInterval,
	}
}

func (c TxConfig) TxReceiptTimeout() time.Duration {
	return c.txReceiptTimeout
}

func (c TxConfig) TxReceiptInterval() time.Duration {
	return c.txReceiptInterval
}

func (c TxConfig) GasLimit() uint64 {
	return c.gasLimit
}