package txclient

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"context"
	"fmt"
	"github.com/lightstreams-network/lightchain/log"
)

var logger = log.NewLogger().With("module", "txclient")

// Dial connects to network and returns directly the official ETH connected Client.
//
// Is recommended to always construct the *ethclient.Client using this Dial method as it only
// returns the Client when the blockchain is fully synced to avoid performing TXs from a non-synced state.
func Dial(gethIpcPath string) (*ethclient.Client, error) {
	logger.Info(fmt.Sprintf("Connecting to ETH client via: '%s'.", gethIpcPath))

	client, err := ethclient.Dial(gethIpcPath)
	if err != nil {
		return nil, err
	}

	syncProgress, err := client.SyncProgress(context.Background())
	if err != nil {
		return nil, err
	}

	if syncProgress != nil {
		return nil, fmt.Errorf("ethereum client is not available till blockchain is fully synced. Sync status: %d/%d", syncProgress.CurrentBlock, syncProgress.HighestBlock)
	}

	return client, nil
}