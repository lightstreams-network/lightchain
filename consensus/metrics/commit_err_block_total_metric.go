package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
)

type CommitErrBlockTotalMetric struct {
	metrics.Counter
}

func (c *CommitErrBlockTotalMetric) Add(delta float64, errorCodeValue string) {
	c.With(errorCodeLabel, errorCodeValue).Add(delta)
}

func NewCommitErrBlockTotalMetric(registry *prometheus.Registry) CommitErrBlockTotalMetric {
	metric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   metricsSubsystem,
		Name:        "commit_err_block_total_counter",
		Help:        "Committed error txs total.",
		ConstLabels: moduleAbciConstLabelValues,
	}, []string{errorCodeLabel})

	registry.Register(metric)

	return CommitErrBlockTotalMetric {
		kitprometheus.NewCounter(metric).With(),
	}
}

func NewNullCommitErrBlockTotalMetric() CommitErrBlockTotalMetric {
	return CommitErrBlockTotalMetric{
		discard.NewCounter(),
	}
}
