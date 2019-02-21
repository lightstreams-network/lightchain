package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type DeliverTxsTotalMetric struct {
	metrics.Counter
}

func NewDeliverTxsTotalMetric(registry *prometheus.Registry) DeliverTxsTotalMetric {
	metric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "deliver_txs_total",
		Help:        "Delivered txs total",
		ConstLabels: moduleAbciConstLabelValues,
	}, []string{})

	registry.Register(metric)

	return DeliverTxsTotalMetric {
		kitprometheus.NewCounter(metric).With(),
	}
}

func NewNullDeliverTxsTotalMetric() DeliverTxsTotalMetric {
	return DeliverTxsTotalMetric{
		discard.NewCounter(),
	}
}
