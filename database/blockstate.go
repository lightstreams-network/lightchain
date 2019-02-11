package database

import (
	"math/big"
	"fmt"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
	tmtCode "github.com/tendermint/tendermint/abci/example/code"
	tmtLog "github.com/tendermint/tendermint/libs/log"
)

// The blockState struct handles processing of TXs included in a block.
//
// It's updated with each ExecuteTx and reset on Persist.
type blockState struct {
	header *types.Header
	parent *types.Block
	state  *state.StateDB

	txIndex      int
	transactions []*types.Transaction
	receipts     types.Receipts
	allLogs      []*types.Log

	totalUsedGas uint64
	gp           *core.GasPool

	logger tmtLog.Logger
}

// Executes TX against the eth blockchain state.
//
// Fetches TX logs and returns new TX receipt. The changes happen only
// inside of the Eth state, not disk!
//
// Logic copied from `core/state_processor.go` `(p *StateProcessor) Process` that gets
// normally executed on block persist.
func (bs *blockState) execTx(bc *core.BlockChain, cfg *eth.Config, chainCfg *params.ChainConfig, tx *types.Transaction) tmtAbciTypes.ResponseDeliverTx {
	snapshot := bs.state.Snapshot()
	bs.state.Prepare(tx.Hash(), common.Hash{}, bs.txIndex)

	receipt, _, err := core.ApplyTransaction(
		chainCfg,
		bc,
		nil, // defaults to address of the author of the header
		bs.gp,
		bs.state,
		bs.header,
		tx,
		&bs.totalUsedGas,
		vm.Config{EnablePreimageRecording: cfg.EnablePreimageRecording},
	)
	if err != nil {
		bs.state.RevertToSnapshot(snapshot)

		return tmtAbciTypes.ResponseDeliverTx{
			Code: tmtCode.CodeTypeUnknownError,
			Log:  fmt.Sprintf("Error applying state TX %v", err),
		}
	}

	logs := bs.state.GetLogs(tx.Hash())

	bs.txIndex++

	bs.transactions = append(bs.transactions, tx)
	bs.receipts = append(bs.receipts, receipt)
	bs.allLogs = append(bs.allLogs, logs...)

	return tmtAbciTypes.ResponseDeliverTx{Code: tmtAbciTypes.CodeTypeOK}
}

// Persist the eth sate, update the header, make a new block and save it to disk.
//
// Returns final application root hash (last block hash).
func (bs *blockState) persist(bc *core.BlockChain, db ethdb.Database) (common.Hash, error) {
	hashArray, err := bs.state.Commit(false)
	if err != nil {
		return common.Hash{}, err
	}
	bs.header.Root = hashArray

	for _, log := range bs.allLogs {
		log.BlockHash = hashArray
	}

	block := types.NewBlock(bs.header, bs.transactions, nil, bs.receipts)
	blockHash := block.Hash()

	// Write block to disk
	_, err = bc.InsertChain([]*types.Block{block})
	if err != nil {
		return common.Hash{}, err
	}

	return blockHash, err
}

func (bs *blockState) updateBlockState(cfg *params.ChainConfig, parentTime uint64, numTx uint64) {
	parentHeader := bs.parent.Header()
	bs.header.Time = new(big.Int).SetUint64(parentTime)
	bs.header.Difficulty = ethash.CalcDifficulty(cfg, parentTime, parentHeader)
	bs.transactions = make([]*types.Transaction, 0, numTx)
	bs.receipts = make([]*types.Receipt, 0, numTx)
	bs.allLogs = make([]*types.Log, 0, numTx)
}