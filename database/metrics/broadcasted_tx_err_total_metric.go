package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type BroadcastedErrTxTotalMetric struct {
	metrics.Counter
}

func (c *BroadcastedErrTxTotalMetric) Add(delta float64, errMsgValue string) {
	c.With(errorMsgLabel, errMsgValue).Add(delta)
}

func NewBroadcastedErrTxTotalMetric(registry *prometheus.Registry) BroadcastedErrTxTotalMetric {
	metric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "broadcasted_err_txs_total",
		Help:        "Broadcasted txs error total.",
	}, []string{errorMsgLabel})

	registry.Register(metric)

	return BroadcastedErrTxTotalMetric {
		kitprometheus.NewCounter(metric).With(),
	}
}

func NewNullBroadcastedErrTxTotalMetric() BroadcastedErrTxTotalMetric {
	return BroadcastedErrTxTotalMetric{
		discard.NewCounter(),
	}
}
