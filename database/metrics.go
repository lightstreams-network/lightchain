package database

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this package.
	MetricsSubsystem = "database"
)

// Metrics contains metrics exposed by this package.
type Metrics struct {
	// ChaindbHeight of the chain.
	ChaindbHeight       metrics.Gauge
	ChaindbNonce        metrics.Gauge
	BroadcastedTxsTotal metrics.Counter
	PersistedTxsTotal   metrics.Counter
	ExecutedTxsTotal    metrics.Counter
}

func PrometheusMetrics(registry *prometheus.Registry, namespace string, labelsAndValues ...string) *Metrics {
	ChaindbHeight := NewGaugeMetric(registry, prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "chaindb_height",
		Help:      "Height of the chaindb.",
	}, labelsAndValues...)

	ChaindbNonce := NewGaugeMetric(registry, prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "chaindb_nonce",
		Help:      "Nonce of the chaindb",
	}, labelsAndValues...)

	BroadcastedTxTotal := NewCounterMetric(registry, prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "broadcasted_txs_total_counter",
		Help:      "Broadcasted txs total.",
	}, labelsAndValues...)

	PersistedTxsTotal := NewCounterMetric(registry, prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "persisted_txs_total_counter",
		Help:      "Persisted txs total",
	}, labelsAndValues...)

	ExecutedTxsTotal := NewCounterMetric(registry, prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "executed_txs_total_counter",
		Help:      "Executed txs total.",
	}, labelsAndValues...)

	return &Metrics{
		ChaindbHeight:       ChaindbHeight,
		ChaindbNonce:        ChaindbNonce,
		BroadcastedTxsTotal: BroadcastedTxTotal,
		PersistedTxsTotal:   PersistedTxsTotal,
		ExecutedTxsTotal:    ExecutedTxsTotal,
	}
}

// NopMetrics returns no-op Metrics.
func NopMetrics() *Metrics {
	return &Metrics{
		ChaindbHeight:       discard.NewGauge(),
		ChaindbNonce:        discard.NewGauge(),
		BroadcastedTxsTotal: discard.NewCounter(),
		PersistedTxsTotal:   discard.NewCounter(),
		ExecutedTxsTotal:    discard.NewCounter(),
	}
}

func NewGaugeMetric(registry *prometheus.Registry, opts prometheus.GaugeOpts, labelsAndValues ...string) metrics.Gauge {
	labels := []string{}
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}

	collection := prometheus.NewGaugeVec(opts, labels)
	registry.MustRegister(collection)
	return kitprometheus.NewGauge(collection).With(labelsAndValues...)
}

func NewCounterMetric(registry *prometheus.Registry, opts prometheus.CounterOpts, labelsAndValues ...string) metrics.Counter {
	labels := []string{}
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}

	collection := prometheus.NewCounterVec(opts, labels)
	registry.MustRegister(collection)
	return kitprometheus.NewCounter(collection).With(labelsAndValues...)
}
