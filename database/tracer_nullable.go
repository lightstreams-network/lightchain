package database

import (
	"github.com/ethereum/go-ethereum/core"
)

var _ tracer = nullEthDBTracer{}

// nullEthDBTracer is used during normal program execution when tracing is disabled.
//
// The point of nullEthDBTracer is to do anything so performance of an app with disabled tracing
// are not affected anyhow.
type nullEthDBTracer struct {
}

func (t nullEthDBTracer) assertPersistedGenesisBlock(genesis core.Genesis) {
}
