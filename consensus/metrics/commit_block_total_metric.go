package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type CommitBlockTotalMetric struct {
	metrics.Counter
}

func NewCommitBlockTotalMetric(registry *prometheus.Registry) CommitBlockTotalMetric {
	metric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "commit_block_total_counter",
		Help:        "Committed txs total.",
		ConstLabels: moduleAbciConstLabelValues,
	}, []string{})

	registry.Register(metric)

	return CommitBlockTotalMetric {
		kitprometheus.NewCounter(metric).With(),
	}
}

func NewNullCommitBlockTotalMetric() CommitBlockTotalMetric {
	return CommitBlockTotalMetric{
		discard.NewCounter(),
	}
}
