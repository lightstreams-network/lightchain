package api

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// A dummy PublicNetAPI in order to overwrite the `ethereum/go-ethereum/internal/ethapi/api.go`.
type PublicNetAPI struct {
	NetworkVersion uint64
}

func NewPublicNetAPI(networkVersion uint64) *PublicNetAPI {
	return &PublicNetAPI{networkVersion}
}

func (n *PublicNetAPI) Listening() bool {
	return true
}

func (n *PublicNetAPI) PeerCount() hexutil.Uint {
	return hexutil.Uint(0)
}

func (n *PublicNetAPI) Version() string {
	return fmt.Sprintf("%d", n.NetworkVersion)
}