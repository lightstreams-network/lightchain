package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type TxsGasTotalMetric struct {
	metrics.Gauge
}

func NewTxsGasTotalMetric(registry *prometheus.Registry) TxsGasTotalMetric {
	metric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "txs_gas_total",
		Help:        "Txs gas total.",
	}, []string{})

	registry.Register(metric)

	return TxsGasTotalMetric {
		kitprometheus.NewGauge(metric).With(),
	}
}

func NewNullTxsGasTotalMetric() TxsGasTotalMetric {
	return TxsGasTotalMetric{
		discard.NewGauge(),
	}
}
