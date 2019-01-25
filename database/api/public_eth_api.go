package api

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/eth"
	"math/big"
	consensusAPI "github.com/lightstreams-network/lightchain/consensus/api"
)

type PublicEthereumAPI struct {
	chainID *big.Int
	e       *eth.Ethereum
	consAPI consensusAPI.API
}

// NewPublicEthereumAPI creates a new Ethereum protocol API for full nodes.
func NewPublicEthereumAPI(chainID *big.Int, e *eth.Ethereum, consAPI consensusAPI.API) *PublicEthereumAPI {
	return &PublicEthereumAPI{chainID, e, consAPI}
}

// Hashrate returns the POW hashrate
func (api *PublicEthereumAPI) Hashrate() hexutil.Uint64 {
	return hexutil.Uint64(0)
}

// ChainId is the EIP-155 replay-protection chain id for the current ethereum chain config.
func (api *PublicEthereumAPI) ChainId() hexutil.Uint64 {
	return (hexutil.Uint64)(api.chainID.Uint64())
}

// Syncing returns whether or not the current node is syncing with other peers. Returns false if not, or a struct
// outlining the state of the sync if it is.
func (e *PublicEthereumAPI) Syncing() (interface{}, error) {
	progress, err := e.consAPI.SyncProgress()
	if err != nil {
		return map[string]interface{} {}, err
	}

	if progress.CurrentBlock == progress.HighestBlock {
		return false, nil
	}

	// This format is necessary to satisfy internal eth API
	// go-ethereum/internal/ethapi/api.go::Syncing()
	return map[string]interface{} {
		"startingBlock": hexutil.Uint64(progress.StartingBlock),
		"currentBlock":  hexutil.Uint64(progress.CurrentBlock),
		"highestBlock":  hexutil.Uint64(progress.HighestBlock),
		"pulledStates":  hexutil.Uint64(progress.PulledStates),
		"knownStates":   hexutil.Uint64(progress.KnownStates),
	}, nil
}