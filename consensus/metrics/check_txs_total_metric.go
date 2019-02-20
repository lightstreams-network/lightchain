package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type CheckTxTotalMetric struct {
	metrics.Counter
}

func NewCheckTxTotalMetric(registry *prometheus.Registry) CheckTxTotalMetric {
	metric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "check_txs_total",
		Help:        "Checked txs total.",
		ConstLabels: moduleAbciConstLabelValues,
	}, []string{})

	registry.Register(metric)

	return CheckTxTotalMetric {
		kitprometheus.NewCounter(metric).With(),
	}
}

func NewNullCheckTxTotalMetric() CheckTxTotalMetric {
	return CheckTxTotalMetric{
		discard.NewCounter(),
	}
}
