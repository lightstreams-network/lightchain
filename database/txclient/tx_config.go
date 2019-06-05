package txclient

import (
	"time"
)

const GasPrice = 500000000000

const TxReceiptInterval = 1 * time.Second
const TxReceiptTimeout = 30 * time.Second

type TxConfig struct {
	gasLimit uint64
	gasPrice uint64
	txReceiptTimeout time.Duration
	txReceiptInterval time.Duration
}

func NewTxConfig(gasLimit uint64, gasPrice uint64, receiptTimeout time.Duration, receiptInterval time.Duration) TxConfig {
	return TxConfig{
		gasLimit,
		gasPrice,
		receiptTimeout,
		receiptInterval,
	}
}

func NewTransferTxConfig() TxConfig {
	return NewTxDefaultConfig(21000)
}

func NewTxDefaultConfig(gasLimit uint64) TxConfig {
	return NewTxConfig(gasLimit, GasPrice, TxReceiptTimeout, TxReceiptInterval)
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

func (c TxConfig) GasPrice() uint64 {
	return c.gasPrice
}