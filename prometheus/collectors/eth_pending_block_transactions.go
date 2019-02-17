package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"context"
	"fmt"
)

type EthPendingBlockTransactions struct {
	ethDialUrl string
	desc   *prometheus.Desc
}

func NewEthPendingBlockTransactions(ethDialUrl string, namespace string) *EthPendingBlockTransactions {
	return &EthPendingBlockTransactions{
		ethDialUrl: ethDialUrl,
		desc: prometheus.NewDesc(
			fmt.Sprintf("%s_%s_eth_pending_block_transactions", namespace, EthereumMetricsSubsystem),
			"the number of transactions in a pending block",
			nil,
			nil,
		),
	}
}

func (collector *EthPendingBlockTransactions) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EthPendingBlockTransactions) Collect(ch chan<- prometheus.Metric) {
	ethClient, err := newEthClient(collector.ethDialUrl)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}
	
	result, err := ethClient.PendingTransactionCount(context.Background())
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, float64(result))
}
