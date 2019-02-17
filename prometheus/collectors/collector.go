package collectors

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	EthereumMetricsSubsystem = "ethereum"
)

type Collectors struct {
	EthPendingBlockTransactions *EthPendingBlockTransactions
	EthSyncing                  *EthSyncing
	EthGenesisBalance           *EthGenesisBalance
}

func EthCollectors(ethDialUrl string, registry *prometheus.Registry, namespace string) *Collectors {
	colEthPendingBlockTransactions := NewEthPendingBlockTransactions(ethDialUrl, namespace)
	colEthSyncing := NewEthSyncing(ethDialUrl, namespace)
	colEthGenesisBalance := NewEthGenesisBalance(ethDialUrl, namespace)

	registry.Register(colEthPendingBlockTransactions)
	registry.Register(colEthSyncing)
	registry.Register(colEthGenesisBalance)

	return &Collectors{
		EthPendingBlockTransactions: colEthPendingBlockTransactions,
		EthSyncing:                  colEthSyncing,
		EthGenesisBalance:           colEthGenesisBalance,
	}
}

// NopMetrics returns no-op 
func EthNullCollectors() *Collectors {
	return &Collectors{}
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
