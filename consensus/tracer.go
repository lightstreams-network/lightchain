package consensus

import (
	"github.com/tendermint/tendermint/config"
	tmtDb "github.com/tendermint/tendermint/libs/db"
	tmtState "github.com/tendermint/tendermint/state"

	"github.com/lightstreams-network/lightchain/network"
	stdtracer "github.com/lightstreams-network/lightchain/tracer"
)

type tracer interface {
	assertPersistedInitStateDb(tmtCfg *config.Config, ntw network.Network)
}

type consensusTracer struct {
	stdtracer.Tracer
}

func newTracer(cfg stdtracer.Config) (tracer, error) {
	trc, err := stdtracer.NewTracer(cfg)
	if err != nil {
		return nil, err
	}
	
	if cfg.ShouldTrace {
		return consensusTracer{trc}, nil		
	}

	return consensusNullTracer{trc}, nil
}

func (trc consensusTracer) assertPersistedInitStateDb(tmtCfg *config.Config, ntw network.Network) {
	trc.Logger.Infow("Tracing whether Tendermint StateDB wrote a valid initial state...", "statedb", tmtCfg.DBPath)

	stateDB := tmtDb.NewDB("state", tmtDb.DBBackendType(tmtCfg.DBBackend), tmtCfg.DBDir())
	defer stateDB.Close()

	persistStateDb := tmtState.LoadState(stateDB)

	protocolVersion, err := ntw.ConsensusProtocolBlockVersion()
	if err != nil {
		trc.Logger.Error(err.Error())
	}
	if persistStateDb.Version.Consensus.Block != protocolVersion {
		trc.Logger.Errorw("Invalid protocol block version persist",
			"expected", protocolVersion,
			"persist", persistStateDb.Version.Consensus.Block)
	}
	
	trc.Logger.Infow("Protocol block version was persisted correctly",
		"block_version", persistStateDb.Version.Consensus.Block)
}
