package database

import (
	"go.uber.org/zap"
	"github.com/ethereum/go-ethereum/core"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
)

type tracerConfig struct {
	shouldTrace bool
	logFilePath string
}

func newTracerConfig(shouldTrace bool, logFilePath string) tracerConfig {
	return tracerConfig{shouldTrace, logFilePath}
}

// tracer is used to trace and assert behaviour of lightchain `database` pkg.
type tracer interface {
	// AssertPersistedGenesisBlock validates if the GenesisBlock was properly saved in disk.
	//
	// Verifies whenever it's possible to re-construct the eth state from disk and
	// asserts accounts balances, nonces, gas settings as defined in Genesis config.
	AssertPersistedGenesisBlock(genesis core.Genesis)
}

// newTracer creates a new tracer instance.
//
// If tracing is disabled, it will return a nullable tracer that doesn't do anything.
//
// DANGER: Tracing is not recommended in production due decrease in performance!
// 		   Use tracing only to debug bugs and for testing purposes.
func newTracer(tracerCfg tracerConfig, chainDataDir string) (tracer, error) {
	if tracerCfg.shouldTrace {
		return newEthDBTracer(chainDataDir, tracerCfg.logFilePath)
	}

	return newNullEthDBTracer(""), nil
}

var _ tracer = ethDBTracer{}

type ethDBTracer struct {
	chainDataDir string
	logger 	     *zap.SugaredLogger
}

func newEthDBTracer(chainDataDir string, logFilePath string) (ethDBTracer, error) {
	logger, err := log.New(logFilePath)
	logger = logger.With("engine", "tracer")
	if err != nil {
		return ethDBTracer{}, err
	}

	return ethDBTracer{
		chainDataDir,
		logger,
	}, nil
}

func (t ethDBTracer) AssertPersistedGenesisBlock(genesis core.Genesis) {
	t.logger.Infow("Tracing if ETH DB wrote a valid genesis block to disk...", "chaindata", t.chainDataDir)

	chainDb, err := ethdb.NewLDBDatabase(t.chainDataDir, 0, 0)
	if err != nil {
		t.logger.Errorw("unable to open LDB db", "err", err)
		return
	}
	defer chainDb.Close()

	head := rawdb.ReadHeadBlockHash(chainDb)
	block := rawdb.ReadBlock(chainDb, head, 0)

	stateDB, err := state.New(block.Root(), state.NewDatabase(chainDb))
	if err != nil {
		t.logger.Errorw("unable to open new stateDB", "err", err)
		return
	}

	for addr, account := range genesis.Alloc {
		persistedBalance := stateDB.GetBalance(addr)
		if persistedBalance.Cmp(account.Balance) != 0 {
			t.logger.Errorw(
				"balance defined in genesis was not properly persisted",
				"acc", addr.Hex(),
				"genesis_balance", account.Balance.String(),
				"persisted_balance", persistedBalance.String(),
			)
			continue
		}

		t.logger.Infow(
			"balance defined in genesis was properly persisted",
			"acc", addr.Hex(),
			"genesis_balance", account.Balance.String(),
			"persisted_balance", persistedBalance.String(),
		)
	}
}