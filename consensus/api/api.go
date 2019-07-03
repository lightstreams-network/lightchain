package api

import (
	eth "github.com/ethereum/go-ethereum"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	tmtcTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtrpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	tmtTypes "github.com/tendermint/tendermint/types"
	tmtCore "github.com/tendermint/tendermint/rpc/core"
	"bytes"
	"fmt"
)

// API is a main consensus interface exposing functionalities to `database` and other packages.
type API interface {
	SyncProgress() (eth.SyncProgress, error)
	BroadcastTx(tx ethTypes.Transaction) error
	Status() (tmtcTypes.ResultStatus, error)
	NetInfo() (tmtcTypes.ResultNetInfo, error)
}

type consensusApi struct {
	isRunning func() bool
	status func(ctx *tmtrpctypes.Context) (*tmtcTypes.ResultStatus, error)
	broadcastTx func (ctx *tmtrpctypes.Context, tx tmtTypes.Tx) (*tmtcTypes.ResultBroadcastTx, error)
	netInfo func (ctx *tmtrpctypes.Context) (*tmtcTypes.ResultNetInfo, error)
}

var _ API = &consensusApi{}

func NewConsensusApi(isRunning func() bool) API {
	return &consensusApi{
		isRunning: isRunning,
		status: tmtCore.Status, 
		broadcastTx: tmtCore.BroadcastTxSync,
		netInfo: tmtCore.NetInfo,
	}
}

func (a *consensusApi) SyncProgress() (eth.SyncProgress, error) {
	status, err := a.status(nil)
	if err != nil {
		return eth.SyncProgress{}, err
	}

	currentBlock := uint64(status.SyncInfo.LatestBlockHeight)
	highestBlock := uint64(status.SyncInfo.LatestBlockHeight)

	if status.SyncInfo.CatchingUp {
		// Latest validator height is unknown because the syncing happens
		// block by block until there are no more and then its considered done
		// if there is no more validators with current state + 1.
		//
		// In the future there could be a precise solution to this if Tendermint
		// would expose this information in the status.
		highestBlock = currentBlock + 1
	}

	return eth.SyncProgress{
		0,
		currentBlock,
		highestBlock,
		0,
		0,
	}, nil
}

func (a *consensusApi) BroadcastTx(tx ethTypes.Transaction) error {
	if !a.isRunning() {
		return fmt.Errorf("Consensus node is not running")
	}

	buf := new(bytes.Buffer)
	if err := tx.EncodeRLP(buf); err != nil {
		return err
	}

	_, err := a.broadcastTx(nil, buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (a *consensusApi) Status() (tmtcTypes.ResultStatus, error) {
	if !a.isRunning() {
		return tmtcTypes.ResultStatus{}, fmt.Errorf("Consensus node is not running")
	}
	
	status, err := tmtCore.Status(nil)
	if err != nil {
		return tmtcTypes.ResultStatus{}, err
	}

	return *status, nil
}

func (a *consensusApi) NetInfo() (tmtcTypes.ResultNetInfo, error) {
	if !a.isRunning() {
		return tmtcTypes.ResultNetInfo{}, fmt.Errorf("Consensus node is not running")
	}
	
	status, err := tmtCore.NetInfo(nil)
	if err != nil {
		return tmtcTypes.ResultNetInfo{}, err
	}

	return *status, nil
}