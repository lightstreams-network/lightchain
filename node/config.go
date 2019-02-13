package node

import (
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/prometheus"
)

type Config struct {
	DataDir string
	consensusCfg consensus.Config
	dbCfg database.Config
	prometheusCfg prometheus.Config
}

func NewConfig(dataDir string,
	consensusCfg consensus.Config,
	dbCfg database.Config,
	prometheusCfg prometheus.Config,
) Config {
	return Config {
		dataDir,
		consensusCfg,
		dbCfg,
		prometheusCfg,
	}
}