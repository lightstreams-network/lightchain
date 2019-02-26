package consensus

import (
	"github.com/tendermint/tendermint/config"

	"github.com/lightstreams-network/lightchain/setup"
	stdtracer "github.com/lightstreams-network/lightchain/tracer"
)

type consensusNullTracer struct {
	stdtracer.Tracer
}

func (trc consensusNullTracer) assertPersistedInitStateDb(tmtCfg *config.Config, ntw setup.Network) {
	return
}
