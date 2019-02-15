package consensus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/lightstreams-network/lightchain/metrics"
)

const (
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this package.
	MetricsSubsystem = "consensus"
)

var (
	ModuleAbciLabelAndValues = []string{"module", "abci"}
)

// Metrics contains metrics exposed by this package.
type Metrics struct {
	CheckTxsTotal       metrics.Counter
	CheckErrTxsTotal    metrics.CounterSet
	DeliverTxsTotal     metrics.Counter
	DeliverErrTxsTotal  metrics.CounterSet
	CommitBlockTotal    metrics.Counter
	CommitErrBlockTotal metrics.CounterSet
}

func TrackedMetrics(registry *prometheus.Registry, namespace string, labelsAndValues ...string) *Metrics {

	checkTxsTotal := metrics.NewCounterMetric(registry, prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "check_txs_total_counter",
		Help:      "Checked txs total",
	}, append(labelsAndValues, ModuleAbciLabelAndValues...)...)

	checkErrTxsTotal := metrics.NewCounterSetMetric(registry, prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "check_err_txs_total_counter",
		Help:      "Checked error txs total.",
	}, []string{"error_code"}, append(labelsAndValues, ModuleAbciLabelAndValues...)...)

	deliverTxsTotal := metrics.NewCounterMetric(registry, prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "deliver_txs_total_counter",
		Help:      "Delivered txs total",
	}, append(labelsAndValues, ModuleAbciLabelAndValues...)...)

	deliverErrTxsTotal := metrics.NewCounterSetMetric(registry, prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "deliver_err_txs_total_counter",
		Help:      "Delivered error txs total.",
	}, []string{"error_code"}, append(labelsAndValues, ModuleAbciLabelAndValues...)...)

	commitBlockTotal := metrics.NewCounterMetric(registry, prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "commit_block_total_counter",
		Help:      "Commited txs total",
	}, append(labelsAndValues, ModuleAbciLabelAndValues...)...)

	commitErrBlockTotal := metrics.NewCounterSetMetric(registry, prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: MetricsSubsystem,
		Name:      "commit_err_block_total_counter",
		Help:      "Commited error txs total.",
	}, []string{"error_code"}, append(labelsAndValues, ModuleAbciLabelAndValues...)...)

	return &Metrics{
		CheckTxsTotal:       checkTxsTotal,
		CheckErrTxsTotal:    checkErrTxsTotal,
		DeliverTxsTotal:     deliverTxsTotal,
		DeliverErrTxsTotal:  deliverErrTxsTotal,
		CommitBlockTotal:    commitBlockTotal,
		CommitErrBlockTotal: commitErrBlockTotal,
	}
}

func TrackedNullMetrics() *Metrics {
	return &Metrics{
		CheckTxsTotal:       metrics.NewCounterDiscard(),
		CheckErrTxsTotal:    metrics.NewCounterSetDiscard(),
		DeliverTxsTotal:     metrics.NewCounterDiscard(),
		DeliverErrTxsTotal:  metrics.NewCounterSetDiscard(),
		CommitBlockTotal:    metrics.NewCounterDiscard(),
		CommitErrBlockTotal: metrics.NewCounterSetDiscard(),
	}
}
