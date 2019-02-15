package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"

	db "github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
	ethMetrics "github.com/lightstreams-network/lightchain/metrics/ethereum"
)

type MetricsProvider struct {
	Database  *db.Metrics
	Consensus *consensus.Metrics
	EthMetrics   *ethMetrics.Metrics
}

func NewMetricsProvider(registry *prometheus.Registry, namespace string, gethIpcPath string) MetricsProvider {
	return MetricsProvider{
		Database:   db.TrackedMetrics(registry, namespace),
		Consensus:  consensus.TrackedMetrics(registry, namespace),
		EthMetrics: ethMetrics.TrackedMetrics(gethIpcPath, registry, namespace),
	}
}

func NewNullMetricsProvider() MetricsProvider {
	return MetricsProvider{
		Database:   db.TrackedNullMetrics(),
		Consensus:  consensus.TrackedNullMetrics(),
		EthMetrics: ethMetrics.TrackedNullMetrics(),
	}
}
