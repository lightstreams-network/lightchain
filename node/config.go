package node

import (
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/prometheus"
	"github.com/lightstreams-network/lightchain/tracer"
)

type Config struct {
	DataDir       string
	consensusCfg  consensus.Config
	dbCfg         database.Config
	prometheusCfg prometheus.Config
	tracerCfg     tracer.Config
}

func NewConfig(dataDir string, consensusCfg consensus.Config, dbCfg database.Config, prometheusCfg prometheus.Config, tracerCfg tracer.Config) Config {
	return Config{
		dataDir,
		consensusCfg,
		dbCfg,
		prometheusCfg,
		tracerCfg,
	}
}