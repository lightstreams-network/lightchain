package consensus

import (
	"fmt"
	"bytes"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/ethereum/go-ethereum/rpc"
	tmtLog "github.com/tendermint/tendermint/libs/log"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
	abciTypes "github.com/lightstreams-network/lightchain/consensus/types"
)

type TendermintABCI struct {
	db           *database.Database
	checkTxState *state.StateDB
	ethRPCClient *rpc.Client
	logger       tmtLog.Logger

	getCurrentDBState func() (*state.StateDB, error)
	getCurrentBlock   func() *ethTypes.Block
}

var _ tmtAbciTypes.Application = &TendermintABCI{}

// NewTendermintABCI creates a new instance of ABCI with clean block's state.
//
// ABCI maintains multiple connections (consensus, info, mempool).
func NewTendermintABCI(db *database.Database, ethRPCClient *rpc.Client) (*TendermintABCI, error) {
	txState, err := db.Ethereum().BlockChain().State()
	if err != nil {
		return nil, err
	}

	abci := &TendermintABCI{
		db:                db,
		ethRPCClient:      ethRPCClient,
		getCurrentDBState: db.Ethereum().BlockChain().State,
		getCurrentBlock:   db.Ethereum().BlockChain().CurrentBlock,
		checkTxState:      txState.Copy(),
		logger:            log.NewLogger().With("engine", "consensus", "module", "ABCI"),
	}

	err = abci.ResetBlockState()
	if err != nil {
		return nil, err
	}

	return abci, nil
}

// InitChain is called upon genesis.
//
// Can be used to define validators set and consensus params on the application side.
//
// Response:
//     - If `ResponseInitChain.Validators` is empty, the initial validator set will be the `RequestInitChain.Validators`
//
func (abci *TendermintABCI) InitChain(req tmtAbciTypes.RequestInitChain) tmtAbciTypes.ResponseInitChain {
	abci.logger.Debug("Initializing chain", "chain_id", req.ChainId)

	return tmtAbciTypes.ResponseInitChain{}
}

// BeginBlock signals the beginning of a new block.
//
// Flow:
//		1. BeginBlock <- ******
//		2. CheckTx
//	    3. ExecuteTx
//		4. EndBlock
//		5. Commit
//		6. CheckTx (clean mempool from TXs not included in committed block)
//
// The header contains the height, timestamp, and more - it exactly matches the Tendermint block header.
// The LastCommitInfo and ByzantineValidators can be used to determine rewards and punishments for the validators.
//
// Response:
//		- Optional Key-Value tags for filtering and indexing
//
func (abci *TendermintABCI) BeginBlock(req tmtAbciTypes.RequestBeginBlock) tmtAbciTypes.ResponseBeginBlock {
	abci.logger.Debug("Beginning new block", "hash", req.Hash)
	abci.db.UpdateHeaderWithTimeInfo(&req.Header)

	return tmtAbciTypes.ResponseBeginBlock{}
}

// ExecuteTx executes the transaction against Ethereum block's work state.
//
// Flow:
//		1. BeginBlock
//		2. CheckTx
//	    3. ExecuteTx <- ******
//		4. EndBlock
//		5. Commit
//		6. CheckTx (clean mempool from TXs not included in committed block)
//
// Tendermint runs CheckTx and ExecuteTx concurrently with each other,
// though on distinct ABCI connections - the mempool connection and the consensus connection, respectively.
//
// Response:
//		- If the transaction is valid, returns CodeType.OK
//		- Keys and values in Tags must be UTF-8 encoded strings
// 		  E.g: ("account.owner": "Bob", "balance": "100.0", "time": "2018-01-02T12:30:00Z")
//
func (abci *TendermintABCI) DeliverTx(txBytes []byte) tmtAbciTypes.ResponseDeliverTx {
	tx, err := decodeRLP(txBytes)
	if err != nil {
		abci.logger.Info("Received invalid transaction", "hash", tx, "err", err)

		return tmtAbciTypes.ResponseDeliverTx{Code: uint32(abciTypes.ErrEncodingError.Code), Log: err.Error()}
	}

	abci.logger.Info("Delivering TX", "hash", tx.Hash().String(), "nonce", tx.Nonce(), "cost", tx.Cost(), "gas", tx.Gas(), "gas_price", tx.GasPrice())

	res := abci.db.ExecuteTx(tx)
	if res.IsErr() {
		abci.logger.Error("Error delivering TX to DB", "hash", tx.Hash().String(), "err", res.Log)

		return res
	}

	abci.logger.Info("TX delivered.", "tx", tx.Hash().String())

	return tmtAbciTypes.ResponseDeliverTx{Code: tmtAbciTypes.CodeTypeOK}
}

// EndBlock signals the end of a block.
//
// Flow:
//		1. BeginBlock
//		2. CheckTx
//	    3. ExecuteTx
//		4. EndBlock <- ******
//		5. Commit
//		6. CheckTx (clean mempool from TXs not included in committed block)
//
// An opportunity to propose changes to a validator set.
//
// Response:
// 		- Validator updates returned for block H
//			- apply to the NextValidatorsHash of block H+1
//			- apply to the ValidatorsHash (and thus the validator set) for block H+2
//			- apply to the RequestBeginBlock.LastCommitInfo (ie. the last validator set) for block H+3
//		- Consensus params returned for block H apply for block H+1
//
func (abci *TendermintABCI) EndBlock(req tmtAbciTypes.RequestEndBlock) tmtAbciTypes.ResponseEndBlock {
	abci.logger.Debug(fmt.Sprintf("Ending new block at height '%d'", req.Height))

	return tmtAbciTypes.ResponseEndBlock{}
}

// Persist persist the application state.
//
// Flow:
//		1. BeginBlock
//		2. CheckTx
//	    3. ExecuteTx
//		4. EndBlock
//		5. Commit <- ******
//		6. CheckTx (clean mempool from TXs not included in committed block)
//
// Response:
//		- Return a Merkle root hash of the application state.
//	      It's critical that all application instances return the same hash.
// 		  If not, they will not be able to agree on the next block,
// 		  because the hash is included in the next block!
//
func (abci *TendermintABCI) Commit() tmtAbciTypes.ResponseCommit {
	rootHash := abci.getCurrentBlock().Root()
	blockHash, err := abci.db.Persist(abci.RewardReceiver())
	if err != nil {
		abci.logger.Error("Error getting latest database state", "err", err)

		return tmtAbciTypes.ResponseCommit{Data: rootHash.Bytes()}
	}

	nextRootHash := abci.getCurrentBlock().Root()
	ethState, err := abci.getCurrentDBState()
	if err != nil {
		abci.logger.Error("Error getting latest state", "err", err)

		return tmtAbciTypes.ResponseCommit{Data: nextRootHash.Bytes()}
	}

	abci.logger.Info("Committing state", "blockHash", blockHash.Hex())

	abci.checkTxState = ethState.Copy()
	if err != nil {
		abci.logger.Error("Error committing latest state", "err", err)

		return tmtAbciTypes.ResponseCommit{Data: nextRootHash.Bytes()}
	}

	return tmtAbciTypes.ResponseCommit{Data: nextRootHash.Bytes()}
}

// RewardReceiver returns the receiving address based on the selected strategy.
func (abci *TendermintABCI) RewardReceiver() common.Address {
	return common.Address{}
}

// ResetBlockState resets the block's processing state.
func (abci *TendermintABCI) ResetBlockState() error {
	return abci.db.ResetBlockState(abci.RewardReceiver())
}

// RLP decode database transaction using go-database impl `go-ethereum/tree/v1.8.11/rlp`.
func decodeRLP(txBytes []byte) (*ethTypes.Transaction, error) {
	tx := new(ethTypes.Transaction)
	rlpStream := rlp.NewStream(bytes.NewBuffer(txBytes), 0)
	if err := tx.DecodeRLP(rlpStream); err != nil {
		return nil, err
	}
	return tx, nil
}