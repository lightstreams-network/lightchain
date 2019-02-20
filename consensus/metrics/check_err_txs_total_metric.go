package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type CheckErrTxsTotalMetric struct {
	metrics.Counter
}

func (c *CheckErrTxsTotalMetric) Add(delta float64, errorCodeValue string) {
	c.With(errorCodeLabel, errorCodeValue).Add(delta)
}

func NewCheckErrTxsTotalMetric(registry *prometheus.Registry) CheckErrTxsTotalMetric {
	metric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "check_err_txs_total",
		Help:        "Checked txs error total.",
		ConstLabels: moduleAbciConstLabelValues,
	}, []string{errorCodeLabel})

	registry.Register(metric)

	return CheckErrTxsTotalMetric{
		kitprometheus.NewCounter(metric).With(),
	}
}

func NewNullCheckErrTrxTotalMetric() CheckErrTxsTotalMetric {
	return CheckErrTxsTotalMetric{
		discard.NewCounter(),
	}
}
