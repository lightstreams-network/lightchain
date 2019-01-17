package api

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// A dummy PublicNetAPI in order to overwrite the `ethereum/go-ethereum/internal/ethapi/api.go`.
type NullPublicNetAPI struct {
	networkVersion uint64
}

func NewNullPublicNetAPI(networkVersion uint64) *NullPublicNetAPI {
	return &NullPublicNetAPI{networkVersion}
}

func (n *NullPublicNetAPI) Listening() bool {
	return true
}

func (n *NullPublicNetAPI) PeerCount() hexutil.Uint {
	return hexutil.Uint(0)
}

func (n *NullPublicNetAPI) Version() string {
	return fmt.Sprintf("%d", n.networkVersion)
}