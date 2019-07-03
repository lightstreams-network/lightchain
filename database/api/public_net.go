package api

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	consensusAPI "github.com/lightstreams-network/lightchain/consensus/api"
)

// A dummy PublicNetAPI in order to overwrite the `ethereum/go-ethereum/internal/ethapi/api.go`.
type PublicNetAPI struct {
	NetworkVersion uint64
	consAPI consensusAPI.API
}

func NewPublicNetAPI(networkVersion uint64, consAPI consensusAPI.API) *PublicNetAPI {
	return &PublicNetAPI{networkVersion, consAPI}
}

func (n *PublicNetAPI) Listening() bool {
	return true
}

func (n *PublicNetAPI) PeerCount() hexutil.Uint {
	netStatus, err := n.consAPI.NetInfo()
	if err != nil {
		return hexutil.Uint(0)		
	}

	return hexutil.Uint(netStatus.NPeers)
}

func (n *PublicNetAPI) Version() string {
	return fmt.Sprintf("%d", n.NetworkVersion)
}