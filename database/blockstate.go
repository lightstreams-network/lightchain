package database

import (
	"math/big"
	"fmt"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	tmtCode "github.com/tendermint/tendermint/abci/example/code"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
)

// The work struct handles block processing.
// It's updated with each DeliverTx and reset on Commit.
type workState struct {
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

//// nolint: unparam
//func (ws *workState) accumulateRewards(strategy *coreUtils.Strategy) {
//	// @Deprecated: Chain Engine is responsible of the rewards
//}

// Runs ApplyTransaction against the eth blockchain, fetches any logs,
// and appends the tx, receipt, and logs.
func (ws *workState) deliverTx(bc *core.BlockChain,
	config *eth.Config,
	chainConfig *params.ChainConfig,
	blockHash common.Hash,
	tx *ethTypes.Transaction) tmtAbciTypes.ResponseDeliverTx {

	ws.state.Prepare(tx.Hash(), blockHash, ws.txIndex)
	receipt, _, err := core.ApplyTransaction(
		chainConfig,
		bc,
		nil, // defaults to address of the author of the header
		ws.gp,
		ws.state,
		ws.header,
		tx,
		&ws.totalUsedGas,
		vm.Config{EnablePreimageRecording: config.EnablePreimageRecording},
	)
	if err != nil {
		return tmtAbciTypes.ResponseDeliverTx{
			Code: tmtCode.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Got %v", err)}
	}

	logs := ws.state.GetLogs(tx.Hash())

	ws.txIndex++

	// The slices are allocated in updateHeaderWithTimeInfo
	ws.transactions = append(ws.transactions, tx)
	ws.receipts = append(ws.receipts, receipt)
	ws.allLogs = append(ws.allLogs, logs...)

	return tmtAbciTypes.ResponseDeliverTx{Code: tmtAbciTypes.CodeTypeOK}
}

// Commit the eth ethState, update the header, make a new block and add it to
// the eth blockchain. The application root hash is the hash of the
// eth block.
func (ws *workState) commit(bc *core.BlockChain, db ethdb.Database) (common.Hash, error) {

	// Commit eth ethState and update the header.
	hashArray, err := ws.state.Commit(false)
	if err != nil {
		return common.Hash{}, err
	}
	ws.header.Root = hashArray

	for _, log := range ws.allLogs {
		log.BlockHash = hashArray
	}

	// Create block object and compute final commit hash (hash of the eth
	// block).
	block := ethTypes.NewBlock(ws.header, ws.transactions, nil, ws.receipts)
	blockHash := block.Hash()

	// Save the block to disk.
	// log.Info("Committing block", "stateHash", hashArray, "blockHash", blockHash)
	_, err = bc.InsertChain([]*ethTypes.Block{block})
	if err != nil {
		// log.Info("Error inserting eth block in chain", "err", err)
		return common.Hash{}, err
	}
	return blockHash, err
}

func (ws *workState) updateHeaderWithTimeInfo(
	config *params.ChainConfig, parentTime uint64, numTx uint64) {

	parentHeader := ws.parent.Header()
	ws.header.Time = new(big.Int).SetUint64(parentTime)
	ws.header.Difficulty = ethash.CalcDifficulty(config, parentTime, parentHeader)
	ws.transactions = make([]*ethTypes.Transaction, 0, numTx)
	ws.receipts = make([]*ethTypes.Receipt, 0, numTx)
	ws.allLogs = make([]*ethTypes.Log, 0, numTx)
}