package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"context"
	"fmt"
)

type EthSyncing struct {
	ethDialUrl string
	desc   *prometheus.Desc
}

func NewEthSyncing(ethDialUrl string, namespace string) *EthSyncing {
	return &EthSyncing{
		ethDialUrl: ethDialUrl,
		desc: prometheus.NewDesc(
			fmt.Sprintf("%s_%s_eth_syncing", namespace, EthereumMetricsSubsystem),
			"the blockchain is syncing",
			nil,
			nil,
		),
	}
}

func (collector *EthSyncing) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EthSyncing) Collect(ch chan<- prometheus.Metric) {
	ethClient, err := newEthClient(collector.ethDialUrl)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
	}

	syncProgress, err := ethClient.SyncProgress(context.Background())
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	if syncProgress == nil {
		// Blockchain is synced
		ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, 1)
	}
}
