package consensus

import (
	"github.com/tendermint/tendermint/config"

	"github.com/lightstreams-network/lightchain/network"
	stdtracer "github.com/lightstreams-network/lightchain/tracer"
)

type consensusNullTracer struct {
	stdtracer.Tracer
}

func (trc consensusNullTracer) assertPersistedInitStateDb(tmtCfg *config.Config, ntw network.Network) {
	return
}
