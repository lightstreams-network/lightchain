package database

import (
	"math/big"
	"sync"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/params"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
)

//----------------------------------------------------------------------
// EthState manages concurrent access to the intermediate workState object
// The eth tx pool fires TxPreEvent in a go-routine,
// and the miner subscribes to this in another go-routine and processes the tx onto
// an intermediate ethState. We used to use `unsafe` to overwrite the miner, but this
// didn't work because it didn't affect the already launched go-routines.
// So instead we introduce the Pending API in a small commit in go-eth
// so we don't even start the miner there, and instead manage the intermediate ethState from here.
// In the same commit we also fire the TxPreEvent synchronously so the order is preserved,
// instead of using a go-routine.

type EthState struct {
	ethereum  *eth.Ethereum
	ethConfig *eth.Config

	mtx  sync.Mutex
	work workState // latest working ethState
}

// After NewEthState, call SetEthereum and SetEthConfig.
func NewEthState() *EthState {
	return &EthState{
		ethereum:  nil, // set with SetEthereum
		ethConfig: nil, // set with SetEthConfig
	}
}

func (es *EthState) SetEthereum(ethereum *eth.Ethereum) {
	es.ethereum = ethereum
}

func (es *EthState) SetEthConfig(ethConfig *eth.Config) {
	es.ethConfig = ethConfig
}

// Execute the transaction.
func (es *EthState) DeliverTx(tx *ethTypes.Transaction) tmtAbciTypes.ResponseDeliverTx {
	es.mtx.Lock()
	defer es.mtx.Unlock()

	bc := es.ethereum.BlockChain()
	chainCfg := es.ethereum.APIBackend.ChainConfig()
	blockHash := common.Hash{}
	return es.work.deliverTx(bc, es.ethConfig, chainCfg, blockHash, tx)
}

//// Accumulate validator rewards.
//func (es *EthState) AccumulateRewards(strategy *coreUtils.Strategy) {
//	es.mtx.Lock()
//	defer es.mtx.Unlock()
//	// @TODO (ggarri): Pending to define reward strategy
//	//ethState.eth.Engine().Finalize(ethState.work.ethState, ethState.work.header, []*ethTypes.Header{})
//	es.work.accumulateRewards(strategy)
//}

// Commit and reset the work.
func (es *EthState) Commit(receiver common.Address) (common.Hash, error) {
	es.mtx.Lock()
	defer es.mtx.Unlock()

	blockHash, err := es.work.commit(es.ethereum.BlockChain(), es.ethereum.ChainDb())
	if err != nil {
		return common.Hash{}, err
	}

	err = es.resetWorkState(receiver)
	if err != nil {
		return common.Hash{}, err
	}

	return blockHash, err
}

func (es *EthState) ResetWorkState(receiver common.Address) error {
	es.mtx.Lock()
	defer es.mtx.Unlock()

	return es.resetWorkState(receiver)
}

func (es *EthState) resetWorkState(receiver common.Address) error {

	blockchain := es.ethereum.BlockChain()
	bcState, err := blockchain.State()
	if err != nil {
		return err
	}

	currentBlock := blockchain.CurrentBlock()
	ethHeader := newBlockHeader(receiver, currentBlock)

	es.work = workState{
		header:       ethHeader,
		parent:       currentBlock,
		state:        bcState,
		txIndex:      0,
		totalUsedGas: uint64(0),
		gp:           new(core.GasPool).AddGas(ethHeader.GasLimit),
	}
	return nil
}

func (es *EthState) UpdateHeaderWithTimeInfo(
	config *params.ChainConfig, parentTime uint64, numTx uint64) {

	es.mtx.Lock()
	defer es.mtx.Unlock()

	es.work.updateHeaderWithTimeInfo(config, parentTime, numTx)
}

func (es *EthState) GasLimit() *core.GasPool {
	return es.work.gp
}

//----------------------------------------------------------------------
// Implements: miner.Pending API (our custom patch to go-eth)

// Return a new block and a copy of the ethState from the latest work.
// #unstable
func (es *EthState) Pending() (*ethTypes.Block, *state.StateDB) {
	es.mtx.Lock()
	defer es.mtx.Unlock()

	return ethTypes.NewBlock(
		es.work.header,
		es.work.transactions,
		nil,
		es.work.receipts,
	), es.work.state.Copy()
}

// Create a new block header from the previous block.
func newBlockHeader(receiver common.Address, prevBlock *ethTypes.Block) *ethTypes.Header {
	return &ethTypes.Header{
		Number:     prevBlock.Number().Add(prevBlock.Number(), big.NewInt(1)),
		ParentHash: prevBlock.Hash(),
		GasLimit:   core.CalcGasLimit(prevBlock, prevBlock.GasLimit(), prevBlock.GasLimit()),
		Coinbase:   receiver,
	}
}
