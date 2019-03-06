package database

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/eth/downloader"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
	dbAPI "github.com/lightstreams-network/lightchain/database/api"
	consensusAPI "github.com/lightstreams-network/lightchain/consensus/api"
	tmtLog "github.com/tendermint/tendermint/libs/log"
	
	"github.com/lightstreams-network/lightchain/database/metrics"
	"github.com/lightstreams-network/lightchain/database/web3"
)

// Database manages the underlying ethereum state for storage and processing
// and maintains the connection to Tendermint for forwarding txs.

// Database handles the chain database and VM.
type Database struct {
	eth    *eth.Ethereum
	ethCfg *eth.Config

	ethTxSub event.Subscription
	ethTxsCh chan core.NewTxsEvent

	ethState *EthState

	consAPI consensusAPI.API
	logger  tmtLog.Logger
	metrics metrics.Metrics
}

func NewDatabase(ctx *node.ServiceContext, ethCfg *eth.Config, consAPI consensusAPI.API, logger tmtLog.Logger, metrics metrics.Metrics) (*Database, error) {
	ethereum, err := eth.New(ctx, ethCfg)
	if err != nil {
		return nil, err
	}

	currentBlock := ethereum.BlockChain().CurrentBlock()
	ethereum.EventMux().Post(core.ChainHeadEvent{currentBlock})

	ethereum.BlockChain().SetValidator(NullBlockValidator{})

	db := &Database{
		eth:      ethereum,
		ethCfg:   ethCfg,
		ethState: NewEthState(ethereum, ethCfg, logger),
		consAPI:  consAPI,
		logger:   logger,
		metrics:  metrics,
	}

	return db, nil
}

func (db *Database) Ethereum() *eth.Ethereum {
	return db.eth
}

func (db *Database) Config() *eth.Config {
	return db.ethCfg
}

// ExecuteTx appends a transaction to the current block.
func (db *Database) ExecuteTx(tx *ethTypes.Transaction) tmtAbciTypes.ResponseDeliverTx {
	db.logger.Info("Executing DB TX", "hash", tx.Hash().Hex(), "nonce", tx.Nonce())

	db.updateExecuteTxMetrics(tx)

	return db.ethState.ExecuteTx(tx)
}

// Persist finalises the current block and writes it to disk.
func (db *Database) Persist(receiver common.Address) (common.Hash, error) {
	db.logger.Info("Persisting DB state", "data", db.ethState.blockState)
	db.metrics.PersistedTxsTotal.Add(float64(len(db.ethState.blockState.transactions)))
	db.metrics.ChaindbHeight.Set(float64(db.ethState.blockState.header.Number.Uint64()))
	return db.ethState.Persist(receiver)
}

// ResetBlockState resets the in-memory block's processing state.
func (db *Database) ResetBlockState(receiver common.Address) error { 
	db.logger.Debug("Resetting ethereum DB state", "receiver", receiver.Hex())
	return db.ethState.ResetBlockState(receiver)
}

// UpdateBlockState uses the tendermint header to update the eth header.
func (db *Database) UpdateBlockState(tmHeader *tmtAbciTypes.Header) {
	db.logger.Debug("Updating ethereum DB state")
	db.ethState.UpdateBlockState(
		db.eth.APIBackend.ChainConfig(),
		uint64(tmHeader.Time.Unix()),
		uint64(tmHeader.GetNumTxs()),
	)
}

// GasLimit returns the maximum gas per block.
func (db *Database) GasLimit() uint64 {
	return db.ethState.GasLimit().Gas()
}

// APIs returns the collection of Ethereum RPC services.
//
// Overwrites go-ethereum/eth/backend.go::APIs().
//
// Some of the API methods must be re-implemented to support Ethereum web3 features
// due to dependency on Tendermint, e.g Syncing().
func (db *Database) APIs() []rpc.API {
	ethAPIs := db.Ethereum().APIs()
	newAPIs := []rpc.API{}

	for _, v := range ethAPIs {
		if isDefaultAPI(v.Namespace) {
			continue
		}

		if _, ok := v.Service.(*eth.PublicMinerAPI); ok {
			continue
		}

		if v.Namespace == "net" {
			v.Service = dbAPI.NewPublicNetAPI(db.ethCfg.NetworkId)
		}

		if _, ok := v.Service.(*eth.PublicEthereumAPI); ok {
			v.Service = dbAPI.NewPublicEthereumAPI(
				db.ethCfg.Genesis.Config.ChainID,
				db.eth,
				db.consAPI,
			)
		}

		if _, ok := v.Service.(downloader.PublicDownloaderAPI); ok {
			v.Service = dbAPI.NewPublicDownloaderAPI()
		}

		newAPIs = append(newAPIs, v)
	}

	return newAPIs
}

func (db *Database) Start(_ *p2p.Server) error {
	go db.txBroadcastLoop()
	return nil
}

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

func (db *Database) Protocols() []p2p.Protocol {
	return nil
}

func (db *Database) updateExecuteTxMetrics(tx *ethTypes.Transaction) {
	db.metrics.ExecutedTxsTotal.Add(1)

	txCost, _ := web3.WeiToPhoton(tx.Cost()).Float64()
	db.metrics.TxsCostTotal.Add(txCost)

	txGasWei := new(big.Int).Mul(big.NewInt(int64(tx.Gas())), tx.GasPrice())
	txGas, _ := web3.WeiToPhoton(txGasWei).Float64()
	db.metrics.TxsGasTotal.Add(txGas)

	db.metrics.TxsSizeTotal.Add(float64(tx.Size()))
}

func isDefaultAPI(namespace string) bool {
	return namespace == "miner" || namespace == "debug" || namespace == "admin"
}
