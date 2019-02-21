package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "lightchain"
	ethereumMetricsSubsystem = "ethereum"
)

type Collectors struct {
	EthPendingBlockTransactions *EthPendingBlockTransactions
	EthSyncing                  *EthSyncing
	EthGenesisBalance           *EthGenesisBalance
}

func NewCollectors(registry *prometheus.Registry, ethDialUrl string) Collectors {
	colEthPendingBlockTransactions := NewEthPendingBlockTransactions(ethDialUrl, namespace)
	colEthSyncing := NewEthSyncing(ethDialUrl, namespace)
	colEthGenesisBalance := NewEthGenesisBalance(ethDialUrl, namespace)

	registry.Register(colEthPendingBlockTransactions)
	registry.Register(colEthSyncing)
	registry.Register(colEthGenesisBalance)

	return Collectors{
		EthPendingBlockTransactions: colEthPendingBlockTransactions,
		EthSyncing:                  colEthSyncing,
		EthGenesisBalance:           colEthGenesisBalance,
	}
}
