package database

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	tmtCode "github.com/tendermint/tendermint/abci/example/code"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
)

// The blockState struct handles processing of TXs included in a block.
//
// It's updated with each ExecuteTx and reset on Persist.
type blockState struct {
	header *ethTypes.Header
	parent *ethTypes.Block
	state  *state.StateDB

	txIndex      int
	transactions []*ethTypes.Transaction
	receipts     ethTypes.Receipts
	allLogs      []*ethTypes.Log

	totalUsedGas uint64
	gp           *core.GasPool
}

// Executes TX against the eth blockchain state.
//
// Fetches TX logs and returns new TX receipt. The changes happen only
// inside of the Eth state, not disk!
//
// Logic copied from `core/state_processor.go` `(p *StateProcessor) Process` that gets
// normally executed on block persist.
func (bs *blockState) execTx(bc *core.BlockChain, config *eth.Config, chainConfig *params.ChainConfig, blockHash common.Hash, tx *ethTypes.Transaction) tmtAbciTypes.ResponseDeliverTx {
	// TODO: Investigate if snapshot should be used `snapshot := bs.state.Snapshot()`
	bs.state.Prepare(tx.Hash(), blockHash, bs.txIndex)
	receipt, _, err := core.ApplyTransaction(
		chainConfig,
		bc,
		nil, // defaults to address of the author of the header
		bs.gp,
		bs.state,
		bs.header,
		tx,
		&bs.totalUsedGas,
		vm.Config{EnablePreimageRecording: config.EnablePreimageRecording},
	)
	if err != nil {
		// TODO: investigate if snapshot should be used `bs.state.RevertToSnapshot(snapshot)`
		return tmtAbciTypes.ResponseDeliverTx{Code: tmtCode.CodeTypeEncodingError, Log: fmt.Sprintf("Error applying state TX %v", err)}
	}

	logs := bs.state.GetLogs(tx.Hash())

	bs.txIndex++

	// The slices are allocated in updateBlockState
	bs.transactions = append(bs.transactions, tx)
	bs.receipts = append(bs.receipts, receipt)
	bs.allLogs = append(bs.allLogs, logs...)

	return tmtAbciTypes.ResponseDeliverTx{Code: tmtAbciTypes.CodeTypeOK}
}

// Persist the eth sate, update the header, make a new block and save it to disk.
//
// Returns final application root hash (last block app state hash).
func (bs *blockState) persist(bc *core.BlockChain, db ethdb.Database) (rootHash common.Hash, err error) {
	rootHash, err = bs.state.Commit(false)
	if err != nil {
		return common.Hash{}, err
	}

	bs.header.Root = rootHash
	for _, log := range bs.allLogs {
		log.BlockHash = bs.header.Root
	}

	// Write block to disk
	block := ethTypes.NewBlock(bs.header, bs.transactions, nil, bs.receipts)
	_, err = bc.InsertChain([]*ethTypes.Block{block})
	if err != nil {
		return common.Hash{}, err
	}

	return rootHash, nil
}

func (bs *blockState) updateBlockState(config *params.ChainConfig, parentTime uint64, numTx uint64) {
	parentHeader := bs.parent.Header()
	bs.header.Time = new(big.Int).SetUint64(parentTime).Uint64()
	bs.header.Difficulty = ethash.CalcDifficulty(config, parentTime, parentHeader)
	bs.transactions = make([]*ethTypes.Transaction, 0, numTx)
	bs.receipts = make([]*ethTypes.Receipt, 0, numTx)
	bs.allLogs = make([]*ethTypes.Log, 0, numTx)
}
