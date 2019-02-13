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
	// Height of the chain.
	Height metrics.Gauge
	NodeTxsTotalCounter metrics.Counter
}

func PrometheusMetrics(registry *prometheus.Registry, namespace string, labelsAndValues ...string) *Metrics {
	HeightMetric := NewGaugeMetric(registry, prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "height",
		Help:      "Height of the chain.",
	}, labelsAndValues...)
	
	NodeTxsTotalCounter := NewCounterMetric(registry, prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "node_txs_total_counter",
		Help:      "Local captured txs.",
	}, labelsAndValues...)

	return &Metrics{
		Height: HeightMetric,
		NodeTxsTotalCounter: NodeTxsTotalCounter,
	}
}

// NopMetrics returns no-op Metrics.
func NopMetrics() *Metrics {
	return &Metrics{
		Height: discard.NewGauge(),
		NodeTxsTotalCounter: discard.NewCounter(),
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
