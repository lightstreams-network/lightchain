package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type ChainDbHeightMetric struct {
	metrics.Gauge
}

func NewChainDbHeightMetric(registry *prometheus.Registry) ChainDbHeightMetric {
	metric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "chaindb_height",
		Help:        "Height of the chaindb.",
	}, []string{})

	registry.Register(metric)

	return ChainDbHeightMetric{
		kitprometheus.NewGauge(metric).With(),
	}
}

func NewNullChainDbHeightMetric() ChainDbHeightMetric {
	return ChainDbHeightMetric{
		discard.NewGauge(),
	}
}
