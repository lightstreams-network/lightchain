package node

import (
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/prometheus"
	"github.com/lightstreams-network/lightchain/tracer"
	"github.com/lightstreams-network/lightchain/governance"
)

type Config struct {
	DataDir       string
	consensusCfg  consensus.Config
	dbCfg         database.Config
	prometheusCfg prometheus.Config
	tracerCfg     tracer.Config
	governanceCfg    governance.Config
}

func NewConfig(
	dataDir string,
	consensusCfg consensus.Config,
	dbCfg database.Config,
	prometheusCfg prometheus.Config,
	tracerCfg tracer.Config,
	governanceCfg governance.Config,
) Config {
	return Config{
		dataDir,
		consensusCfg,
		dbCfg,
		prometheusCfg,
		tracerCfg,
		governanceCfg,
	}
}

func (c Config) DbCfg() database.Config {
	return c.dbCfg
}

func (c Config) TracerCfg() tracer.Config {
	return c.tracerCfg
}

func (c Config) GovernanceCfg() governance.Config {
	return c.governanceCfg
}