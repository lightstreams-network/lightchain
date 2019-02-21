package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace        = "lightchain"
	metricsSubsystem = "database"
)

const errorMsgLabel = "error_msg"

type Metrics struct {
	ChaindbHeight          ChainDbHeightMetric
	BroadcastedTxsTotal    BroadcastedTxTotalMetric
	BroadcastedErrTxsTotal BroadcastedErrTxTotalMetric
	PersistedTxsTotal      PersistedTxsTotalMetric
	ExecutedTxsTotal       ExecutedTxsTotalMetric
	TxsSizeTotal           TxsSizeTotalMetric
	TxsCostTotal           TxsCostTotalMetric
	TxsGasTotal            TxsGasTotalMetric
}

func NewMetrics(registry *prometheus.Registry) Metrics {
	return Metrics{
		ChaindbHeight:          NewChainDbHeightMetric(registry),
		BroadcastedTxsTotal:    NewBroadcastedTxTotalMetric(registry),
		BroadcastedErrTxsTotal: NewBroadcastedErrTxTotalMetric(registry),
		PersistedTxsTotal:      NewPersistedTxsTotalMetric(registry),
		ExecutedTxsTotal:       NewExecutedTxsTotalMetric(registry),
		TxsSizeTotal:           NewTxsSizeTotalMetric(registry),
		TxsCostTotal:           NewTxsCostTotalMetric(registry),
		TxsGasTotal:            NewTxsGasTotalMetric(registry),
	}
}

// NopMetrics returns no-op metrics.
func NewNullMetrics() Metrics {
	return Metrics{
		ChaindbHeight:          NewNullChainDbHeightMetric(),
		BroadcastedTxsTotal:    NewNullBroadcastedTxTotalMetric(),
		BroadcastedErrTxsTotal: NewNullBroadcastedErrTxTotalMetric(),
		PersistedTxsTotal:      NewNullPersistedTxsTotalMetric(),
		ExecutedTxsTotal:       NewNullExecutedTxsTotalMetric(),
		TxsSizeTotal:           NewNullTxsSizeTotalMetric(),
		TxsCostTotal:           NewNullTxsCostTotalMetric(),
		TxsGasTotal:            NewNullTxsGasTotalMetric(),
	}
}
