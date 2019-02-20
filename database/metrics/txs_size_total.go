package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type TxsSizeTotalMetric struct {
	metrics.Gauge
}

func NewTxsSizeTotalMetric(registry *prometheus.Registry) TxsSizeTotalMetric {
	metric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "txs_size_total",
		Help:        "Txs size total.",
	}, []string{})

	registry.Register(metric)

	return TxsSizeTotalMetric {
		kitprometheus.NewGauge(metric).With(),
	}
}

func NewNullTxsSizeTotalMetric() TxsSizeTotalMetric {
	return TxsSizeTotalMetric{
		discard.NewGauge(),
	}
}
