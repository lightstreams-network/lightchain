package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type DeliverErrTxsTotalMetric struct {
	metrics.Counter
}

func (c *DeliverErrTxsTotalMetric) Add(delta float64, errorCodeValue string) {
	c.With(errorCodeLabel, errorCodeValue).Add(delta)
}

func NewDeliverErrTxsTotalMetric(registry *prometheus.Registry) DeliverErrTxsTotalMetric {
	metric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "deliver_err_txs_total",
		Help:        "Delivered error txs total.",
		ConstLabels: moduleAbciConstLabelValues,
	}, []string{errorCodeLabel})

	registry.Register(metric)

	return DeliverErrTxsTotalMetric {
		kitprometheus.NewCounter(metric).With(),
	}
}

func NewNullDeliverErrTxsTotalMetric() DeliverErrTxsTotalMetric {
	return DeliverErrTxsTotalMetric{
		discard.NewCounter(),
	}
}
