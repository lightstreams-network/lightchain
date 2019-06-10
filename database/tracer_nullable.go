package database

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	tmtConfig "github.com/tendermint/tendermint/config"
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

func (t nullEthDBTracer) AssertPersistedValidatorSetContract(contractAddress common.Address, ownerAddress common.Address) {
}

func (t nullEthDBTracer) AssertPersistedValidatorSetAddValidator(tmtCfg tmtConfig.Config, validatorPubKey string, rewardedAddress common.Address) {
}