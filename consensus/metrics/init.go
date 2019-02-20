package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "lightchain"
	metricsSubsystem = "consensus"
)

var (
	moduleAbciConstLabelValues = map[string]string{"module": "abci"}
)

const errorCodeLabel = "error_code"

// Metrics contains metrics exposed by this package.
type Metrics struct {
	CheckTxsTotal       CheckTxTotalMetric
	CheckErrTxsTotal    CheckErrTxsTotalMetric
	DeliverTxsTotal     DeliverTxsTotalMetric
	DeliverErrTxsTotal  DeliverErrTxsTotalMetric
	CommitBlockTotal    CommitBlockTotalMetric
	CommitErrBlockTotal CommitErrBlockTotalMetric
}

func NewMetrics(registry *prometheus.Registry) Metrics {
	return Metrics{
		CheckTxsTotal:       NewCheckTxTotalMetric(registry),
		CheckErrTxsTotal:    NewCheckErrTxsTotalMetric(registry),
		DeliverTxsTotal:     NewDeliverTxsTotalMetric(registry),
		DeliverErrTxsTotal:  NewDeliverErrTxsTotalMetric(registry),
		CommitBlockTotal:    NewCommitBlockTotalMetric(registry),
		CommitErrBlockTotal: NewCommitErrBlockTotalMetric(registry),
	}
}

func NewNullMetrics() Metrics {
	return Metrics{
		CheckTxsTotal:       NewNullCheckTxTotalMetric(),
		CheckErrTxsTotal:    NewNullCheckErrTrxTotalMetric(),
		DeliverTxsTotal:     NewNullDeliverTxsTotalMetric(),
		DeliverErrTxsTotal:  NewNullDeliverErrTxsTotalMetric(),
		CommitBlockTotal:    NewNullCommitBlockTotalMetric(),
		CommitErrBlockTotal: NewNullCommitErrBlockTotalMetric(),
	}
}
