package database

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var _ Tracer = nullEthDBTracer{}

// nullEthDBTracer is used during normal program execution when tracing is disabled.
//
// The point of nullEthDBTracer is to do anything so performance of an app with disabled tracing
// are not affected anyhow.
type nullEthDBTracer struct {
}

func (t nullEthDBTracer) AssertPersistedGenesisBlock(genesis core.Genesis) {
}

func (t nullEthDBTracer) AssertPostTxSimulationState(from common.Address, tx *types.Transaction) {
}