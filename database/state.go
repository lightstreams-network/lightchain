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
	tmtLog "github.com/tendermint/tendermint/libs/log"
)

//----------------------------------------------------------------------
// EthState manages concurrent access to the intermediate blockState object
// The eth tx pool fires TxPreEvent in a go-routine,
// and the miner subscribes to this in another go-routine and processes the tx onto
// an intermediate ethState. We used to use `unsafe` to overwrite the miner, but this
// didn't work because it didn't affect the already launched go-routines.
// So instead we introduce the Pending API in a small persist in go-eth
// so we don't even start the miner there, and instead manage the intermediate ethState from here.
// In the same persist we also fire the TxPreEvent synchronously so the order is preserved,
// instead of using a go-routine.
//
// TODO: Rewrite this description, it doesn't make sense and this "EthState" Struct is even more suspicious.
type EthState struct {
	ethereum  *eth.Ethereum
	ethConfig *eth.Config

	mtx        sync.Mutex
	blockState blockState

	logger tmtLog.Logger
}

func NewEthState(ethereum *eth.Ethereum, ethCfg *eth.Config, logger tmtLog.Logger) *EthState {
	return &EthState{
		ethereum:  ethereum,
		ethConfig: ethCfg,
		logger: logger,
	}
}

// Executes TX against the eth blockchain state.
//
// The changes happen only inside of the Eth state, not disk!
func (es *EthState) ExecuteTx(tx *ethTypes.Transaction) tmtAbciTypes.ResponseDeliverTx {
	es.mtx.Lock()
	defer es.mtx.Unlock()

	bc := es.ethereum.BlockChain()
	chainCfg := es.ethereum.APIBackend.ChainConfig()
	blockHash := common.Hash{}
	return es.blockState.execTx(bc, es.ethConfig, chainCfg, blockHash, tx)
}

// Persist the application state to disk.
//
// Persist is called after all TXs are executed against the application state.
// Triggered by ABCI::Commit(), to persist changes introduced in latest block.
//
// Returns the persisted Block.
func (es *EthState) Persist() (ethTypes.Block, error) {
	es.mtx.Lock()
	defer es.mtx.Unlock()

	block, err := es.blockState.persist(es.ethereum.BlockChain(), es.ethereum.ChainDb())
	if err != nil {
		return ethTypes.Block{}, err
	}

	err = es.resetBlockState()
	if err != nil {
		return ethTypes.Block{}, err
	}

	return block, err
}

func (es *EthState) ResetBlockState(receiver common.Address) error {
	es.mtx.Lock()
	defer es.mtx.Unlock()

	return es.resetBlockState()
}

func (es *EthState) resetBlockState() error {
	blockchain := es.ethereum.BlockChain()
	bcState, err := blockchain.State()
	if err != nil {
		return err
	}

	currentBlock := blockchain.CurrentBlock()
	ethHeader := newBlockHeader(common.Address{}, currentBlock)

	es.blockState = blockState{
		header:       ethHeader,
		parent:       currentBlock,
		state:        bcState,
		txIndex:      0,
		totalUsedGas: uint64(0),
		gp:           new(core.GasPool).AddGas(ethHeader.GasLimit),
	}
	return nil
}

func (es *EthState) UpdateBlockState(config *params.ChainConfig, parentTime uint64, numTx uint64, receiver common.Address) {
	es.mtx.Lock()
	defer es.mtx.Unlock()

	es.blockState.updateBlockState(config, parentTime, numTx, receiver)
}

func (es *EthState) GasLimit() *core.GasPool {
	return es.blockState.gp
}

// Implements: miner.Pending API (our custom patch to go-eth).
//
// Return a new block and a copy of the ethState from the latest blockState.
func (es *EthState) Pending() (*ethTypes.Block, *state.StateDB) {
	es.mtx.Lock()
	defer es.mtx.Unlock()

	return ethTypes.NewBlock(
		es.blockState.header,
		es.blockState.transactions,
		nil,
		es.blockState.receipts,
	), es.blockState.state.Copy()
}

func newBlockHeader(receiver common.Address, prevBlock *ethTypes.Block) *ethTypes.Header {
	return &ethTypes.Header{
		Number:     prevBlock.Number().Add(prevBlock.Number(), big.NewInt(1)),
		ParentHash: prevBlock.Hash(),
		GasLimit:   core.CalcGasLimit(prevBlock, prevBlock.GasLimit(), prevBlock.GasLimit()),
		Coinbase:   receiver,
	}
}
