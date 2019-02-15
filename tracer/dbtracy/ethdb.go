package dbtracy

import (
	"go.uber.org/zap"
	"github.com/ethereum/go-ethereum/core"
	"github.com/lightstreams-network/lightchain/log"
)

// Tracer is used to trace and assert behaviour of lightchain `database` pkg.
type Tracer interface {
	// AssertPersistedGenesisBlock validates if the GenesisBlock was properly saved in disk.
	//
	// Verifies whenever it's possible to re-construct the eth state from disk and
	// asserts accounts balances, nonces, gas settings as defined in Genesis config.
	AssertPersistedGenesisBlock(genesis core.Genesis)
}

// New creates a new tracer instance.
//
// If tracing is disabled, it will return a nullable tracer that doesn't do anything.
//
// DANGER: Tracing is not recommended in production due decrease in performance!
// 		   Use tracing only to debug bugs and for testing purposes.
func New(shouldTrace bool, chainDataDir string, logFilePath string) (Tracer, error) {
	if shouldTrace {
		return NewEthDBTracer(chainDataDir, logFilePath)
	}

	return NewNullEthDBTracer(""), nil
}

var _ Tracer = EthDBTracer{}

type EthDBTracer struct {
	chainDataDir string
	logger 	     *zap.SugaredLogger
}

func NewEthDBTracer(chainDataDir string, logFilePath string) (EthDBTracer, error) {
	logger, err := log.New(logFilePath)
	logger = logger.With("engine", "tracer")
	if err != nil {
		return EthDBTracer{}, err
	}

	return EthDBTracer{
		chainDataDir,
		logger,
	}, nil
}