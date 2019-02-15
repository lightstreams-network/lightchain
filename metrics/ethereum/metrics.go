package collector

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	EthereumMetricsSubsystem = "ethereum"
)

type Metrics struct {
	EthPendingBlockTransactions *EthPendingBlockTransactions
	EthSyncing                  *EthSyncing
}

func TrackedMetrics(ethDialUrl string, registry *prometheus.Registry, namespace string) *Metrics {
	colEthPendingBlockTransactions := NewEthPendingBlockTransactions(ethDialUrl, namespace)
	colEthSyncing := NewEthSyncing(ethDialUrl, namespace)

	registry.Register(colEthPendingBlockTransactions)
	registry.Register(colEthSyncing)

	return &Metrics{
		EthPendingBlockTransactions: colEthPendingBlockTransactions,
		EthSyncing:                  colEthSyncing,
	}
}

// NopMetrics returns no-op 
func TrackedNullMetrics() *Metrics {
	return &Metrics{}
}

// Dial connects to network and returns directly the official ETH connected Client.
func newEthClient(conn string) (*ethclient.Client, error) {
	log.Printf("Connecting to ETH client via: '%s'.", conn)

	client, err := ethclient.Dial(conn)
	if err != nil {
		return nil, err
	}

	return client, nil
}
