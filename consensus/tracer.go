package consensus

import (
	"github.com/tendermint/tendermint/config"
	tmtDb "github.com/tendermint/tendermint/libs/db"
	tmtState "github.com/tendermint/tendermint/state"

	"github.com/lightstreams-network/lightchain/setup"
	stdtracer "github.com/lightstreams-network/lightchain/tracer"
)

type tracer interface {
	assertPersistedInitStateDb(tmtCfg *config.Config, ntw setup.Network)
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

func (trc consensusTracer) assertPersistedInitStateDb(tmtCfg *config.Config, ntw setup.Network) {
	trc.Logger.Infow("Tracing whether Tendermint StateDB wrote a valid initial state...", "statedb", tmtCfg.DBPath)

	stateDB := tmtDb.NewDB("state", tmtDb.DBBackendType(tmtCfg.DBBackend), tmtCfg.DBDir())
	defer stateDB.Close()

	persistStateDb := tmtState.LoadState(stateDB)

	switch ntw {
	case setup.SiriusNetwork:
		if persistStateDb.Version.Consensus.Block != siriusProtocolBlockVersion {
			trc.Logger.Errorw("Invalid protocol block version persist",
				"expected", siriusProtocolBlockVersion,
				"persist", persistStateDb.Version.Consensus.Block)
		}
		break;
	case setup.StandaloneNetwork:
		if persistStateDb.Version.Consensus.Block != standaloneProtocolBlockVersion {
			trc.Logger.Errorw("Invalid protocol block version persist",
				"expected", standaloneProtocolBlockVersion,
				"persist", persistStateDb.Version.Consensus.Block)
		}
		break;
	default:
		trc.Logger.Errorw("Invalid network selected", "network", ntw)
		return
	}

	trc.Logger.Infow("Protocol block version was persisted correctly",
		"block_version", persistStateDb.Version.Consensus.Block)
}
