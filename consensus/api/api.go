package api

import (
	eth "github.com/ethereum/go-ethereum"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	tmtTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtRpcClient "github.com/tendermint/tendermint/rpc/lib/client"
	"fmt"
	"bytes"
	"github.com/lightstreams-network/lightchain/log"
	"net/http"
	"io/ioutil"
	"strings"
	"encoding/json"
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

	highestBlock, err := fetchLatestValidatorBlockHeight()
	if err != nil {
		return eth.SyncProgress{}, err
	}
	
	return eth.SyncProgress{
		0,
		uint64(status.SyncInfo.LatestBlockHeight),
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

// The Explorer can be considered as a source of truth of the network
// as connected validator could be out of sync as well.
func fetchLatestValidatorBlockHeight() (uint64, error) {
	url := "https://explorer.lightstreams.io/web3relay"
	payload := strings.NewReader("{\n\t\"action\": \"blockrate\"\n}")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	web3RelayResp := struct {
		BlockHeight uint64 `json:"blockHeight"`
	}{}
	err = json.Unmarshal(resBody, &web3RelayResp)
	if err != nil {
		return 0, err
	}

	return web3RelayResp.BlockHeight, nil
}