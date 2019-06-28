package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type ReplacedBlockTimeTotalMetric struct {
	metrics.Counter
}

func NewReplacedBlockTimeTotalMetric(registry *prometheus.Registry) ReplacedBlockTimeTotalMetric {
	metric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "replaced_block_time_total_counter",
		Help:        "Replaced consensus block time total.",
		ConstLabels: moduleAbciConstLabelValues,
	}, []string{})

	registry.Register(metric)

	return ReplacedBlockTimeTotalMetric{
		kitprometheus.NewCounter(metric).With(),
	}
}

func NewNullReplacedBlockTimeTotalMetric() ReplacedBlockTimeTotalMetric {
	return ReplacedBlockTimeTotalMetric{
		discard.NewCounter(),
	}
}
