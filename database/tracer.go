package database

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"

	stdtracer "github.com/lightstreams-network/lightchain/tracer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"bytes"
	"math/big"
	"github.com/lightstreams-network/lightchain/database/web3"
)

// Tracer is used to trace and assert behaviour of lightchain `database` pkg.
type Tracer interface {
	// AssertPersistedGenesisBlock validates if the GenesisBlock was properly saved in disk.
	//
	// Verifies whenever it's possible to re-construct the eth state from disk and
	// asserts accounts balances, nonces, gas settings as defined in Genesis config.
	AssertPersistedGenesisBlock(genesis core.Genesis)

	// AssertPostTxSimulationState validates if the state saved in disk after TX simulation is correct.
	AssertPostTxSimulationState(from common.Address, tx *types.Transaction)
}

var _ Tracer = EthDBTracer{}

type EthDBTracer struct {
	stdtracer.Tracer
	chainDataDir string
}

// NewTracer creates a new Tracer instance.
//
// If tracing is disabled, it will return a nullable Tracer that doesn't do anything.
//
// DANGER: Tracing is not recommended in production due decrease in performance!
// 		   Use tracing only to debug bugs and for testing purposes.
func NewTracer(cfg stdtracer.Config, chainDataDir string) (Tracer, error) {
	trc, err := stdtracer.NewTracer(cfg)
	if err != nil {
		return nil, err
	}

	if cfg.ShouldTrace {
		return EthDBTracer{trc, chainDataDir}, nil
	}

	return nullEthDBTracer{}, nil
}

func (t EthDBTracer) AssertPersistedGenesisBlock(genesis core.Genesis) {
	t.Logger.Infow("Tracing if ETH DB wrote a valid genesis block to disk...", "chaindata", t.chainDataDir)

	chainDb, err := ethdb.NewLDBDatabase(t.chainDataDir, 0, 0)
	if err != nil {
		t.Logger.Errorw("unable to open LDB db", "err", err)
		return
	}
	defer chainDb.Close()

	head := rawdb.ReadHeadBlockHash(chainDb)
	block := rawdb.ReadBlock(chainDb, head, 0)

	stateDB, err := state.New(block.Root(), state.NewDatabase(chainDb))
	if err != nil {
		t.Logger.Errorw("unable to open new stateDB", "err", err)
		return
	}

	for addr, account := range genesis.Alloc {
		persistedBalance := stateDB.GetBalance(addr)
		if persistedBalance.Cmp(account.Balance) != 0 {
			t.Logger.Errorw(
				"balance defined in genesis was not properly persisted",
				"acc", addr.Hex(),
				"genesis_balance", account.Balance.String(),
				"persisted_balance", persistedBalance.String(),
			)
			continue
		}

		t.Logger.Infow(
			"balance defined in genesis was properly persisted",
			"acc", addr.Hex(),
			"genesis_balance", account.Balance.String(),
			"persisted_balance", persistedBalance.String(),
		)
	}
}

func (t EthDBTracer) AssertPostTxSimulationState(from common.Address, tx *types.Transaction) {
	t.Logger.Infow("Tracing if ETH DB is in a valid state after simulation...", "chaindata", t.chainDataDir)

	chainDb, err := ethdb.NewLDBDatabase(t.chainDataDir, 0, 0)
	if err != nil {
		t.Logger.Errorw("unable to open LDB db", "err", err)
		return
	}
	defer chainDb.Close()

	headHash := rawdb.ReadHeadBlockHash(chainDb)
	headNumber := rawdb.ReadHeaderNumber(chainDb, headHash)
	block := rawdb.ReadBlock(chainDb, headHash, *headNumber)

	stateDB, err := state.New(block.Root(), state.NewDatabase(chainDb))
	if err != nil {
		t.Logger.Errorw("unable to open new stateDB", "err", err)
		return
	}

	// Given the same TX inputs, hash must not change
	expectedTxHash := common.HexToHash("0x11aef34de0533989f78f43439913f5f11cd31d4cbcbd43e52630dd8a1f2d69e8")
	if !bytes.Equal(block.TxHash().Bytes(), expectedTxHash.Bytes()) {
		t.Logger.Errorw("incorrect simulated TX hash", "expected", expectedTxHash.Hex(), "actual", block.TxHash().Hex())
	} else {
		t.Logger.Infow("correct simulated TX hash", "hash", block.TxHash().Hex())
	}

	// Given the same TX changes, block root hash must not change
	expectedRootHash := common.HexToHash("0x61bf842528ab995fc819b306b42ddb31b92aee0328f1cabf4f2cf4d7d23e9c45")
	if !bytes.Equal(block.Root().Bytes(), expectedRootHash.Bytes()) {
		t.Logger.Errorw("incorrect root hash", "expected", expectedRootHash.Hex(), "actual", block.Root().Hex())
	} else {
		t.Logger.Infow("correct root hash", "hash", block.Root().Hex())
	}

	// Coinbase should be 0x0
	expectedCoinbase := common.HexToAddress("0x0")
	if !bytes.Equal(block.Coinbase().Bytes(), expectedCoinbase.Bytes()) {
		t.Logger.Errorw("incorrect coinbase", "expected", expectedCoinbase.Hex(), "actual", block.Coinbase().Hex())
	} else {
		t.Logger.Infow("correct coinbase", "acc", block.Coinbase().Hex())
	}

	// Parent hash is always the genesis block
	expectedParentHash := common.HexToHash("0x55e06fc7b51b31efb053f128068be8b09c86569895d98591ea5790b683770c58")
	if !bytes.Equal(block.ParentHash().Bytes(), expectedParentHash.Bytes()) {
		t.Logger.Errorw("incorrect parent hash", "expected", expectedParentHash.Hex(), "actual", block.ParentHash().Hex())
	} else {
		t.Logger.Infow("correct parent hash", "hash", block.ParentHash().Hex())
	}

	genesisFromBalance, _ := web3.ParseWei("300000000000000000000000000")
	genesisFromBalanceCopy, _ := new(big.Int).SetString(genesisFromBalance.String(), 10)
	expectedFromBalance := new(big.Int).Sub(genesisFromBalanceCopy, tx.Cost())

	fromBalance := stateDB.GetBalance(from)
	if fromBalance.Cmp(expectedFromBalance) != 0 {
		t.Logger.Errorw(
			"incorrect post TX balance",
			"acc", from.Hex(),
			"genesis_balance", genesisFromBalance.String(),
			"expected_balance", expectedFromBalance.String(),
			"post_simulation_balance", fromBalance.String(),
		)
	} else {
		t.Logger.Infow(
			"sender balance after TX simulation is correct",
			"acc", from.Hex(),
			"genesis_balance", genesisFromBalance.String(),
			"expected_balance", expectedFromBalance.String(),
			"post_simulation_balance", fromBalance.String(),
		)
	}

	requiredGasPrice := big.NewInt(MinGasPrice)
	if requiredGasPrice.Cmp(tx.GasPrice()) > 0 {
		t.Logger.Errorw(
			"incorrect gas price. Expected default gas price",
			"required_gas_price", requiredGasPrice.String(),
			"tx_gas_price", tx.GasPrice().String(),
		)
	} else {
		t.Logger.Infow(
			"TX gas price set to default gas price as expected",
			"required_gas_price", requiredGasPrice.String(),
			"tx_gas_price", tx.GasPrice().String(),
		)
	}

	fromNonce := stateDB.GetNonce(from)
	expectedNonce := uint64(1)
	if fromNonce != expectedNonce {
		t.Logger.Errorw(
			"incorrect sender nonce",
			"expected_nonce", expectedNonce,
			"actual_nonce", fromNonce,
		)
	} else {
		t.Logger.Infow(
			"correct sender nonce",
			"expected_nonce", expectedNonce,
			"actual_nonce", fromNonce,
		)
	}
}