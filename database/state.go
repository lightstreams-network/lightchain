package database

import (
	"math/big"
	"sync"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/params"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
	tmtLog "github.com/tendermint/tendermint/libs/log"
)

//----------------------------------------------------------------------
// State manages concurrent access to the intermediate blockState object
// The eth tx pool fires TxPreEvent in a go-routine,
// and the miner subscribes to this in another go-routine and processes the tx onto
// an intermediate state. We used to use `unsafe` to overwrite the miner, but this
// didn't bs because it didn't affect the already launched go-routines.
// So instead we introduce the Pending API in a small persist in go-eth
// so we don't even start the miner there, and instead manage the intermediate state from here.
// In the same persist we also fire the TxPreEvent synchronously so the order is preserved,
// instead of using a go-routine.
//
// TODO: Rewrite this description, it doesn't make sense and this "State" Struct is even more suspicious.

type State struct {
	ethereum  *eth.Ethereum
	ethConfig *eth.Config

	mtx        sync.Mutex
	blockState blockState

	logger tmtLog.Logger
}

func NewState(ethereum *eth.Ethereum, ethCfg *eth.Config, logger tmtLog.Logger) *State {
	return &State{
		ethereum:  ethereum,
		ethConfig: ethCfg,
		logger: logger,
	}
}

// Executes TX against the eth blockchain state.
//
// The changes happen only inside of the Eth state, not disk!
func (s *State) ExecuteTx(tx *ethTypes.Transaction) tmtAbciTypes.ResponseDeliverTx {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	bc := s.ethereum.BlockChain()
	chainCfg := s.ethereum.APIBackend.ChainConfig()

	return s.blockState.execTx(bc, s.ethConfig, chainCfg, tx)
}

// Persist the application state to disk.
//
// Persist is called after all TXs are executed against the application state.
// Triggered by ABCI::Commit(), to persist changes introduced in latest block.
func (s *State) Persist(coinbase common.Address) (common.Hash, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	blockHash, err := s.blockState.persist(s.ethereum.BlockChain(), s.ethereum.ChainDb())
	if err != nil {
		return common.Hash{}, err
	}

	err = s.resetBlockState(coinbase)
	if err != nil {
		return common.Hash{}, err
	}

	return blockHash, err
}

func (s *State) ResetBlockState(coinbase common.Address) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.resetBlockState(coinbase)
}

func (s *State) UpdateBlockState(config *params.ChainConfig, parentTime uint64, numTx uint64) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.blockState.updateBlockState(config, parentTime, numTx)
}

func (s *State) GasLimit() *core.GasPool {
	return s.blockState.gp
}

func (s *State) resetBlockState(coinbase common.Address) error {
	blockchain := s.ethereum.BlockChain()
	bcState, err := blockchain.State()
	if err != nil {
		return err
	}

	currentBlock := blockchain.CurrentBlock()
	ethHeader := newBlockHeader(coinbase, currentBlock)

	s.blockState = blockState{
		header:       ethHeader,
		parent:       currentBlock,
		state:        bcState,
		txIndex:      0,
		totalUsedGas: uint64(0),
		gp:           new(core.GasPool).AddGas(ethHeader.GasLimit),
		logger:       s.logger,
	}

	return nil
}

func newBlockHeader(coinbase common.Address, prevBlock *ethTypes.Block) *ethTypes.Header {
	return &ethTypes.Header{
		Number:     prevBlock.Number().Add(prevBlock.Number(), big.NewInt(1)),
		ParentHash: prevBlock.Hash(),
		GasLimit:   core.CalcGasLimit(prevBlock, prevBlock.GasLimit(), prevBlock.GasLimit()),
		Coinbase:   coinbase,
	}
}