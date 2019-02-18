package database

import (
	"go.uber.org/zap"
	"github.com/ethereum/go-ethereum/core"
)

var _ tracer = nullEthDBTracer{}

// nullEthDBTracer is used during normal program execution when tracing is disabled.
//
// The point of nullEthDBTracer is to do anything so performance of an app with disabled tracing
// are not affected anyhow.
type nullEthDBTracer struct {
	chainDataDir string
	logger 	     *zap.SugaredLogger
}

func newNullEthDBTracer(chainDataDir string) nullEthDBTracer {
	return nullEthDBTracer{}
}

func (t nullEthDBTracer) AssertPersistedGenesisBlock(genesis core.Genesis) {
}