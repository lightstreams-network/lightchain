package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type TxsCostTotalMetric struct {
	metrics.Gauge
}

func NewTxsCostTotalMetric(registry *prometheus.Registry) TxsCostTotalMetric {
	metric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "txs_cost_total",
		Help:        "Txs cost total.",
	}, []string{})

	registry.Register(metric)

	return TxsCostTotalMetric {
		kitprometheus.NewGauge(metric).With(),
	}
}

func NewNullTxsCostTotalMetric() TxsCostTotalMetric {
	return TxsCostTotalMetric{
		discard.NewGauge(),
	}
}
