package database

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/log"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
	tmtRpcClient "github.com/tendermint/tendermint/rpc/lib/client"
)

// Database manages the underlying ethereum state for storage and processing
// and maintains the connection to Tendermint for forwarding txs

// Database handles the chain database and VM
type Database struct {
	eth    *eth.Ethereum
	ethCfg *eth.Config

	ethTxSub event.Subscription
	ethTxsCh chan core.NewTxsEvent

	ethState *EthState

	tmtRPCClient tmtRpcClient.HTTPClient
}

// New creates a new database
func New(ctx *node.ServiceContext, ethCfg *eth.Config, tmtRPCClient tmtRpcClient.HTTPClient) (*Database, error) {
	state := NewEthState()

	// eth.NewNode takes a ServiceContext for the EventMux, the AccountManager,
	// and some basic functions around the DataDir.
	ethereum, err := eth.New(ctx, ethCfg)
	if err != nil {
		return nil, err
	}

	state.SetEthereum(ethereum)
	state.SetEthConfig(ethCfg)

	// send special event to go-eth to switch homestead=true.
	currentBlock := ethereum.BlockChain().CurrentBlock()
	ethereum.EventMux().Post(core.ChainHeadEvent{currentBlock})

	// We don't need PoW/Uncle validation.
	ethereum.BlockChain().SetValidator(NullBlockProcessor{})

	db := &Database{
		eth:          ethereum,
		ethCfg:       ethCfg,
		ethState:     state,
		tmtRPCClient: tmtRPCClient,
	}

	return db, nil
}

// Ethereum returns the underlying the ethereum object.
// #stable
func (db *Database) Ethereum() *eth.Ethereum {
	return db.eth
}

// Config returns the eth.Config.
// #stable
func (db *Database) Config() *eth.Config {
	return db.ethCfg
}

//----------------------------------------------------------------------
// Handle block processing

// DeliverTx appends a transaction to the current block
func (db *Database) DeliverTx(tx *ethTypes.Transaction) tmtAbciTypes.ResponseDeliverTx {
	log.Info("Lightchain.Database::DeliverTx", "tx", tx.Hash().String())
	return db.ethState.DeliverTx(tx)
}

//// AccumulateRewards accumulates the rewards based on the given strategy
//func (db *Database) AccumulateRewards(strategy *coreUtils.Strategy) {
//	db.ethState.AccumulateRewards(strategy)
//}

// Commit finalises the current block
func (db *Database) Commit(receiver common.Address) (common.Hash, error) {
	log.Info("Lightchain.Database::Commit", "data", db.ethState.work)
	return db.ethState.Commit(receiver)
}

// InitEthState initializes the EthState
func (db *Database) InitEthState(receiver common.Address) error {
	log.Info("Lightchain.Database::InitEthState", "data", receiver)
	return db.ethState.ResetWorkState(receiver)
}

// UpdateHeaderWithTimeInfo uses the tendermint header to update the eth header
func (db *Database) UpdateHeaderWithTimeInfo(tmHeader *tmtAbciTypes.Header) {
	db.ethState.UpdateHeaderWithTimeInfo(db.eth.APIBackend.ChainConfig(),
		uint64(tmHeader.Time.Unix()), uint64(tmHeader.GetNumTxs()))
}

// GasLimit returns the maximum gas per block
func (db *Database) GasLimit() uint64 {
	return db.ethState.GasLimit().Gas()
}

//----------------------------------------------------------------------
// Implements: eth.Service

// APIs returns the collection of RPC services the eth package offers.
func (db *Database) APIs() []rpc.API {
	apis := db.Ethereum().APIs()
	retApis := []rpc.API{}
	for _, v := range apis {
		if v.Namespace == "net" {
			v.Service = NewNetRPCService(db.ethCfg.NetworkId)
		}
		if v.Namespace == "miner" {
			continue
		}
		if _, ok := v.Service.(*eth.PublicMinerAPI); ok {
			continue
		}
		retApis = append(retApis, v)
	}
	return retApis
}

// Start implements node.Service, starting all internal goroutines needed by the Ethereum
// protocol implementation.
func (db *Database) Start(_ *p2p.Server) error {
	go db.txBroadcastLoop()
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the Ethereum protocol.
func (db *Database) Stop() error {
	db.ethTxSub.Unsubscribe()
	db.eth.BlockChain().Stop()
	db.eth.Engine().Close()
	db.eth.TxPool().Stop()
	db.eth.ChainDb().Close()
	//if err := db.eth.Stop(); err != nil {
	//	return err
	//}
	return nil
}

// Protocols implements node.Service, returning all the currently configured network protocols to start.
func (db *Database) Protocols() []p2p.Protocol {
	return nil
}

func (db *Database) FlushStateTrieDb() error {
	bc := db.Ethereum().BlockChain()
	triedb := bc.StateCache().TrieDB()
	root := bc.CurrentBlock().Root()
	if err := triedb.Commit(root, true); err != nil {
		return err
	}
	return nil
}
