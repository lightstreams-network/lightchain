package api

import (
	eth "github.com/ethereum/go-ethereum"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	tmtTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtRpcClient "github.com/tendermint/tendermint/rpc/lib/client"
	"fmt"
	"bytes"
	"github.com/lightstreams-network/lightchain/log"
)

// API is a main consensus interface exposing functionalities to `database` and other packages.
type API interface {
	SyncProgress() (eth.SyncProgress, error)
	BroadcastTx(tx ethTypes.Transaction) error
	Status() (tmtTypes.ResultStatus, error)
}

type rpcApi struct {
	client tmtRpcClient.HTTPClient
	logger log.Logger
}

var _ API = &rpcApi{}

func NewRPCApi(rpcListenPort uint) API {
	tendermintLAddr := fmt.Sprintf("tcp://127.0.0.1:%d", rpcListenPort)
	client := tmtRpcClient.NewURIClient(tendermintLAddr)
	tmtTypes.RegisterAmino(client.Codec())


	logger := log.NewLogger()
	logger.With("module", "consensus_api")

	return &rpcApi{client, logger}
}

func (a *rpcApi) SyncProgress() (eth.SyncProgress, error) {
	status := tmtTypes.ResultStatus{}
	_, err := a.client.Call("status", map[string]interface{} {}, &status)
	if err != nil {
		return eth.SyncProgress{}, err
	}

	// Arbitrary values representing not synced state
	currentBlock := uint64(0)
	highestBlock := uint64(1000000)

	if !status.SyncInfo.CatchingUp {
		currentBlock = highestBlock
	}
	
	return eth.SyncProgress{
		0,
		currentBlock,
		highestBlock,
		0,
		0,
	}, nil
}

func (a *rpcApi) BroadcastTx(tx ethTypes.Transaction) error {
	a.logger.Info("Broadcasting tx to Tendermint core via RPC API...")

	syncInfo := new(tmtTypes.SyncInfo)
	buf := new(bytes.Buffer)
	if err := tx.EncodeRLP(buf); err != nil {
		return err
	}

	params := map[string]interface{} {
		"tx": buf.Bytes(),
	}

	_, err := a.client.Call("broadcast_tx_sync", params, syncInfo)
	if err != nil {
		return err
	}

	return nil
}

func (a *rpcApi) Status() (tmtTypes.ResultStatus, error) {
	status := new(tmtTypes.ResultStatus)
	_, err := a.client.Call("status", map[string]interface{}{}, &status)
	if err != nil {
		return tmtTypes.ResultStatus{}, err
	}

	return *status, nil
}