package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/lightstreams-network/lightchain/consensus"
	db "github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/prometheus/collectors"
)

type MetricsProvider struct {
	Database   *db.Metrics
	Consensus  *consensus.Metrics
	EthMetrics *collectors.Collectors
}

func NewMetricsProvider(registry *prometheus.Registry, namespace string, gethIpcPath string) MetricsProvider {
	return MetricsProvider{
		Database:   db.TrackedMetrics(registry, namespace),
		Consensus:  consensus.TrackedMetrics(registry, namespace),
		EthMetrics: collectors.EthCollectors(gethIpcPath, registry, namespace),
	}
}

func NewNullMetricsProvider() MetricsProvider {
	return MetricsProvider{
		Database:   db.TrackedNullMetrics(),
		Consensus:  consensus.TrackedNullMetrics(),
		EthMetrics: collectors.EthNullCollectors(),
	}
}
