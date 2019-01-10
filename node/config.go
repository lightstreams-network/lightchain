package node

import (
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
)

type Config struct {
	DataDir string
	consensusCfg consensus.Config
	dbCfg database.Config
}

func NewConfig(dataDir string, consensusCfg consensus.Config, dbCfg database.Config) Config {
	return Config {
		dataDir,
		consensusCfg,
		dbCfg,
	}
}
