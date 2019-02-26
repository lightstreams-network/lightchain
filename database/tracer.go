package database

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	
	stdtracer "github.com/lightstreams-network/lightchain/tracer"
)

// tracer is used to trace and assert behaviour of lightchain `database` pkg.
type tracer interface {
	// assertPersistedGenesisBlock validates if the GenesisBlock was properly saved in disk.
	//
	// Verifies whenever it's possible to re-construct the eth state from disk and
	// asserts accounts balances, nonces, gas settings as defined in Genesis config.
	assertPersistedGenesisBlock(genesis core.Genesis)
}

var _ tracer = ethDBTracer{}

type ethDBTracer struct {
	stdtracer.Tracer
	chainDataDir string
}

// newTracer creates a new tracer instance.
//
// If tracing is disabled, it will return a nullable tracer that doesn't do anything.
//
// DANGER: Tracing is not recommended in production due decrease in performance!
// 		   Use tracing only to debug bugs and for testing purposes.
func newTracer(cfg stdtracer.Config, chainDataDir string) (tracer, error) {
	trc, err := stdtracer.NewTracer(cfg)
	if err != nil {
		return nil, err
	}
	
	if cfg.ShouldTrace {
		return ethDBTracer{trc, ChainDbPath}, nil 
	}
	
	return nullEthDBTracer{}, nil
}


func (t ethDBTracer) assertPersistedGenesisBlock(genesis core.Genesis) {
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
