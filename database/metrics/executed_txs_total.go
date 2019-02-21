package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type ExecutedTxsTotalMetric struct {
	metrics.Counter
}

func NewExecutedTxsTotalMetric(registry *prometheus.Registry) ExecutedTxsTotalMetric {
	metric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "executed_txs_total",
		Help:        "Executed txs total.",
	}, []string{})

	registry.Register(metric)

	return ExecutedTxsTotalMetric {
		kitprometheus.NewCounter(metric).With(),
	}
}

func NewNullExecutedTxsTotalMetric() ExecutedTxsTotalMetric {
	return ExecutedTxsTotalMetric{
		discard.NewCounter(),
	}
}
