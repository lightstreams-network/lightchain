package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type PersistedTxsTotalMetric struct {
	metrics.Counter
}

func NewPersistedTxsTotalMetric(registry *prometheus.Registry) PersistedTxsTotalMetric {
	metric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "persisted_txs_total",
		Help:        "Persisted txs total.",
	}, []string{})

	registry.Register(metric)

	return PersistedTxsTotalMetric {
		kitprometheus.NewCounter(metric).With(),
	}
}

func NewNullPersistedTxsTotalMetric() PersistedTxsTotalMetric {
	return PersistedTxsTotalMetric{
		discard.NewCounter(),
	}
}
