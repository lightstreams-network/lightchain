package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
	tmtRpcClient "github.com/tendermint/tendermint/rpc/lib/client"

	"github.com/ethereum/go-ethereum/log"
	coreUtils "github.com/lightstreams-network/lightchain/utils"
)

//----------------------------------------------------------------------
// Backend manages the underlying ethereum state for storage and processing,
// and maintains the connection to Tendermint for forwarding txs

// Backend handles the chain database and VM
type Backend struct {
	// backing ethereum structures
	ethereum  *eth.Ethereum
	ethConfig *eth.Config

	// txBroadcastLoop subscription
	txSub event.Subscription
	txsCh chan core.NewTxsEvent

	// EthState
	ethState *EthState

	// client for forwarding txs to Tendermint
	client tmtRpcClient.HTTPClient
}

// NewBackend creates a new Backend
func NewBackend(ctx *node.ServiceContext, ethConfig *eth.Config, client tmtRpcClient.HTTPClient) (*Backend, error) {

	// Create working ethereum state.
	ethState := NewEthState()

	// eth.New takes a ServiceContext for the EventMux, the AccountManager,
	// and some basic functions around the DataDir.
	ethereum, err := eth.New(ctx, ethConfig)
	if err != nil {
		return nil, err
	}

	ethState.SetEthereum(ethereum)
	ethState.SetEthConfig(ethConfig)

	// send special event to go-ethereum to switch homestead=true.
	currentBlock := ethereum.BlockChain().CurrentBlock()
	ethereum.EventMux().Post(core.ChainHeadEvent{currentBlock}) // nolint: vet, errcheck

	// We don't need PoW/Uncle validation.
	ethereum.BlockChain().SetValidator(NullBlockProcessor{})

	ethBackend := &Backend{
		ethereum:  ethereum,
		ethConfig: ethConfig,
		ethState:  ethState,
		client:    client,
	}

	return ethBackend, nil
}

// Ethereum returns the underlying the ethereum object.
// #stable
func (b *Backend) Ethereum() *eth.Ethereum {
	return b.ethereum
}

// Config returns the eth.Config.
// #stable
func (b *Backend) Config() *eth.Config {
	return b.ethConfig
}

//----------------------------------------------------------------------
// Handle block processing

// DeliverTx appends a transaction to the current block
func (b *Backend) DeliverTx(tx *ethTypes.Transaction) tmtAbciTypes.ResponseDeliverTx {
	log.Info("Lightchain.Backend::DeliverTx", "tx", tx.Hash().String())
	return b.ethState.DeliverTx(tx)
}

// AccumulateRewards accumulates the rewards based on the given strategy
func (b *Backend) AccumulateRewards(strategy *coreUtils.Strategy) {
	b.ethState.AccumulateRewards(strategy)
}

// Commit finalises the current block
func (b *Backend) Commit(receiver common.Address) (common.Hash, error) {
	log.Info("Lightchain.Backend::Commit", "data", b.ethState.work)
	return b.ethState.Commit(receiver)
}

// InitEthState initializes the EthState
func (b *Backend) InitEthState(receiver common.Address) error {
	log.Info("Lightchain.Backend::InitEthState", "data", receiver)
	return b.ethState.ResetWorkState(receiver)
}

// UpdateHeaderWithTimeInfo uses the tendermint header to update the ethereum header
func (b *Backend) UpdateHeaderWithTimeInfo(tmHeader *tmtAbciTypes.Header) {
	b.ethState.UpdateHeaderWithTimeInfo(b.ethereum.APIBackend.ChainConfig(),
		uint64(tmHeader.Time.Unix()), uint64(tmHeader.GetNumTxs()))
}

// GasLimit returns the maximum gas per block
func (b *Backend) GasLimit() uint64 {
	return b.ethState.GasLimit().Gas()
}

//----------------------------------------------------------------------
// Implements: node.Service

// APIs returns the collection of RPC services the ethereum package offers.
func (b *Backend) APIs() []rpc.API {
	apis := b.Ethereum().APIs()
	retApis := []rpc.API{}
	for _, v := range apis {
		if v.Namespace == "net" {
			v.Service = NewNetRPCService(b.ethConfig.NetworkId)
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
func (b *Backend) Start(_ *p2p.Server) error {
	go b.txBroadcastLoop()
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the Ethereum protocol.
func (b *Backend) Stop() error {
	b.txSub.Unsubscribe()
	b.ethereum.BlockChain().Stop()
	b.ethereum.Engine().Close()
	b.ethereum.TxPool().Stop()
	b.ethereum.ChainDb().Close()
	//if err := b.ethereum.Stop(); err != nil {
	//	return err
	//}
	return nil
}

// Protocols implements node.Service, returning all the currently configured network protocols to start.
func (b *Backend) Protocols() []p2p.Protocol {
	return nil
}

func (b *Backend) FlushStateTrieDb() error {
	bc := b.Ethereum().BlockChain()
	triedb := bc.StateCache().TrieDB()
	root := bc.CurrentBlock().Root()
	if err := triedb.Commit(root, true); err != nil {
		return err
	}
	return nil
}

//----------------------------------------------------------------------
// We need a block processor that just ignores PoW and uncles and so on

// NullBlockProcessor does not validate anything
// #unstable
type NullBlockProcessor struct{}

// ValidateBody does not validate anything
// #unstable
func (NullBlockProcessor) ValidateBody(*ethTypes.Block) error { return nil }

// ValidateState does not validate anything
// #unstable
func (NullBlockProcessor) ValidateState(block, parent *ethTypes.Block, state *state.StateDB,
	receipts ethTypes.Receipts, usedGas uint64) error {
	return nil
}
