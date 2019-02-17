package database

import (
	
	"github.com/prometheus/client_golang/prometheus"
	"github.com/lightstreams-network/lightchain/prometheus/metrics"
)

const (
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this package.
	MetricsSubsystem = "database"
)

type Metrics struct {
	ChaindbHeight       metrics.Gauge
	BroadcastedTxsTotal metrics.Counter
	PersistedTxsTotal   metrics.Counter
	ExecutedTxsTotal    metrics.Counter
	TxsSizeTotal        metrics.Gauge
	TxsCostTotal        metrics.Gauge
	TxsGasTotal         metrics.Gauge
}

func TrackedMetrics(registry *prometheus.Registry, namespace string, labelsAndValues ...string) *Metrics {
	ChaindbHeight := metrics.NewGaugeMetric(registry, prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "chaindb_height",
		Help:      "Height of the chaindb.",
	}, labelsAndValues...)

	BroadcastedTxTotal := metrics.NewCounterMetric(registry, prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "broadcasted_txs_total_counter",
		Help:      "Broadcasted txs total.",
	}, labelsAndValues...)

	PersistedTxsTotal := metrics.NewCounterMetric(registry, prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "persisted_txs_total_counter",
		Help:      "Persisted txs total.",
	}, labelsAndValues...)

	ExecutedTxsTotal := metrics.NewCounterMetric(registry, prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "executed_txs_total_counter",
		Help:      "Executed txs total.",
	}, labelsAndValues...)
	
	txsSizeTotal := metrics.NewGaugeMetric(registry, prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "txs_size_total",
		Help:      "Txs size total.",
	}, labelsAndValues...)

	txsCostTotal := metrics.NewGaugeMetric(registry, prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "txs_cost_total",
		Help:      "Txs cost total.",
	}, labelsAndValues...)

	txsGasTotal := metrics.NewGaugeMetric(registry, prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "txs_gas_total",
		Help:      "Txs gas total.",
	}, labelsAndValues...)

	return &Metrics{
		ChaindbHeight:       ChaindbHeight,
		BroadcastedTxsTotal: BroadcastedTxTotal,
		PersistedTxsTotal:   PersistedTxsTotal,
		ExecutedTxsTotal:    ExecutedTxsTotal,
		TxsSizeTotal:       txsSizeTotal,
		TxsCostTotal:       txsCostTotal,
		TxsGasTotal:        txsGasTotal,
	}
}

// NopMetrics returns no-op metrics.
func TrackedNullMetrics() *Metrics {
	return &Metrics{
		ChaindbHeight:       metrics.NewGaugeDiscard(),
		BroadcastedTxsTotal: metrics.NewCounterDiscard(),
		PersistedTxsTotal:   metrics.NewCounterDiscard(),
		ExecutedTxsTotal:    metrics.NewCounterDiscard(),
		TxsSizeTotal:        metrics.NewGaugeDiscard(),
		TxsCostTotal:        metrics.NewGaugeDiscard(),
		TxsGasTotal:         metrics.NewGaugeDiscard(),
	}
}

