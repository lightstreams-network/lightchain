package abci

import (
	"fmt"
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/ethereum/go-ethereum/core/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	abciTypes "github.com/lightstreams-network/lightchain/consensus/types"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
	tmtLog "github.com/tendermint/tendermint/libs/log"
)

type ConsensusConnection struct {
	db *database.Database

	replaceTxState    func(*state.StateDB)
	getCurrentDBState func() (*state.StateDB, error)
	getCurrentBlock   func() *ethTypes.Block

	logger tmtLog.Logger
}

func newConsensusConnection(db *database.Database, replaceTxState func(*state.StateDB), logger tmtLog.Logger) (*ConsensusConnection, error) {
	return &ConsensusConnection{
		db,
		replaceTxState,
		db.Ethereum().BlockChain().State,
		db.Ethereum().BlockChain().CurrentBlock,
		logger,
	}, nil
}

// InitChain is called upon genesis.
//
// Can be used to define validators set and consensus params on the application side.
//
// Response:
//     - If `ResponseInitChain.Validators` is empty, the initial validator set will be the `RequestInitChain.Validators`
func (cc *ConsensusConnection) InitChain(req tmtAbciTypes.RequestInitChain) tmtAbciTypes.ResponseInitChain {
	cc.logger.Debug("Initializing chain", "chain_id", req.ChainId)

	return tmtAbciTypes.ResponseInitChain{}
}

// BeginBlock signals the beginning of a new block.
//
// The header contains the height, timestamp, and more - it exactly matches the Tendermint block header.
// The `req.LastCommitInfo` and `req.ByzantineValidators` attributes can be used to determine rewards and punishments
// for the validators.
//
// Response:
//		- Optional Key-Value tags for filtering and indexing
func (cc *ConsensusConnection) BeginBlock(req tmtAbciTypes.RequestBeginBlock) tmtAbciTypes.ResponseBeginBlock {
	cc.logger.Debug("Beginning new block", "hash", req.Hash)
	cc.db.UpdateHeaderWithTimeInfo(&req.Header)

	return tmtAbciTypes.ResponseBeginBlock{}
}

// DeliverTx executes the transaction against Ethereum block's work state.
//
// Response:
//		- If the transaction is valid, returns CodeType.OK
//		- Keys and values in Tags must be UTF-8 encoded strings
// 		  E.g: ("account.owner": "Bob", "balance": "100.0", "time": "2018-01-02T12:30:00Z")
func (cc *ConsensusConnection) DeliverTx(txBytes []byte) tmtAbciTypes.ResponseDeliverTx {
	tx, err := decodeRLP(txBytes)
	if err != nil {
		cc.logger.Info("Received invalid transaction", "hash", tx, "err", err)

		return tmtAbciTypes.ResponseDeliverTx{Code: uint32(abciTypes.ErrEncodingError.Code), Log: err.Error()}
	}

	cc.logger.Info("Delivering TX", "hash", tx.Hash().String(), "nonce", tx.Nonce(), "cost", tx.Cost(), "gas", tx.Gas(), "gas_price", tx.GasPrice())

	res := cc.db.ExecuteTx(tx)
	if res.IsErr() {
		cc.logger.Error("Error delivering TX to DB", "hash", tx.Hash().String(), "err", res.Log)

		return res
	}

	cc.logger.Info("TX delivered.", "tx", tx.Hash().String())

	return tmtAbciTypes.ResponseDeliverTx{Code: tmtAbciTypes.CodeTypeOK}
}

// EndBlock signals the end of a block.
//
// An opportunity to propose changes to a validator set.
//
// Response:
// 		- Validator updates returned for block H
//			- apply to the NextValidatorsHash of block H+1
//			- apply to the ValidatorsHash (and thus the validator set) for block H+2
//			- apply to the RequestBeginBlock.LastCommitInfo (ie. the last validator set) for block H+3
//		- Consensus params returned for block H apply for block H+1
func (cc *ConsensusConnection) EndBlock(req tmtAbciTypes.RequestEndBlock) tmtAbciTypes.ResponseEndBlock {
	cc.logger.Debug(fmt.Sprintf("Ending new block at height '%d'", req.Height))

	return tmtAbciTypes.ResponseEndBlock{}
}

// Persist persist the application state.
//
// Response:
//		- Return a Merkle root hash of the application state.
//	      It's critical that all application instances return the same hash.
// 		  If not, they will not be able to agree on the next block,
// 		  because the hash is included in the next block!
func (cc *ConsensusConnection) Commit() tmtAbciTypes.ResponseCommit {
	rootHash := cc.getCurrentBlock().Root()
	blockHash, err := cc.db.Persist(cc.RewardReceiver())
	if err != nil {
		cc.logger.Error("Error getting latest database state", "err", err)

		return tmtAbciTypes.ResponseCommit{Data: rootHash.Bytes()}
	}

	nextRootHash := cc.getCurrentBlock().Root()
	ethState, err := cc.getCurrentDBState()
	if err != nil {
		cc.logger.Error("Error getting latest state", "err", err)

		return tmtAbciTypes.ResponseCommit{Data: nextRootHash.Bytes()}
	}

	cc.logger.Info("Committing state", "blockHash", blockHash.Hex())

	nextRootHash, err = ethState.Commit(false)

	if err != nil {
		cc.logger.Error("Error committing latest state", "err", err)

		return tmtAbciTypes.ResponseCommit{Data: nextRootHash.Bytes()}
	}

	cc.replaceTxState(ethState.Copy())

	return tmtAbciTypes.ResponseCommit{Data: nextRootHash.Bytes()}
}

// RewardReceiver returns the receiving address based on the selected strategy.
func (cc *ConsensusConnection) RewardReceiver() common.Address {
	return common.Address{}
}

// ResetBlockState resets the block's processing state.
func (cc *ConsensusConnection) ResetBlockState() error {
	return cc.db.ResetBlockState(cc.RewardReceiver())
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