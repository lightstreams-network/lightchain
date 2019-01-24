package api

import (
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/rpc"
	"context"
)

// PublicDownloaderAPI overwrites the `ethereum/go-ethereum/eth/downloader/api.go`
// and disables its features.
type PublicDownloaderAPI struct {
}

func NewPublicDownloaderAPI() *PublicDownloaderAPI {
	return &PublicDownloaderAPI{}
}

func (api *PublicDownloaderAPI) Syncing(ctx context.Context) (*rpc.Subscription, error) {
	return nil, featureNotSupportedErr()
}

func (api *PublicDownloaderAPI) SubscribeSyncStatus(status chan interface{}) *downloader.SyncStatusSubscription {
	return nil
}