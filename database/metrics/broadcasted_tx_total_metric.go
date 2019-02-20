package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type BroadcastedTxTotalMetric struct {
	metrics.Counter
}

func NewBroadcastedTxTotalMetric(registry *prometheus.Registry) BroadcastedTxTotalMetric {
	metric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "broadcasted_txs_total",
		Help:        "Broadcasted txs total.",
	}, []string{})

	registry.Register(metric)

	return BroadcastedTxTotalMetric {
		kitprometheus.NewCounter(metric).With(),
	}
}

func NewNullBroadcastedTxTotalMetric() BroadcastedTxTotalMetric {
	return BroadcastedTxTotalMetric{
		discard.NewCounter(),
	}
}
