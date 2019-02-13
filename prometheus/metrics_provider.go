package prometheus

import (
	db "github.com/lightstreams-network/lightchain/database"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsProvider struct {
	Database *db.Metrics	
}

func NewMetricsProvider(registry *prometheus.Registry, namespace string) MetricsProvider {
	return MetricsProvider {
		Database: db.PrometheusMetrics(registry, namespace),
	}
}

func NewNopMetricsProvider() MetricsProvider {
	return MetricsProvider{
		Database: db.NopMetrics(),
	}
}
