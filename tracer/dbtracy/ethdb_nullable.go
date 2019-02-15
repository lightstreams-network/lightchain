package dbtracy

import (
	"go.uber.org/zap"
	"github.com/ethereum/go-ethereum/core"
)

var _ Tracer = NullEthDBTracer{}

// NullEthDBTracer is used during normal program execution when tracing is disabled.
//
// The point of NullEthDBTracer is to do anything so performance of an app with disabled tracing
// are not affected anyhow.
type NullEthDBTracer struct {
	chainDataDir string
	logger 	     *zap.SugaredLogger
}

func NewNullEthDBTracer(chainDataDir string) NullEthDBTracer {
	return NullEthDBTracer{}
}

func (t NullEthDBTracer) AssertPersistedGenesisBlock(genesis core.Genesis) {
}