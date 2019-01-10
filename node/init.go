package node

import (
	"github.com/lightstreams-network/lightchain/utils"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/database"
)

var logger = log.NewLogger()

func InitNode(cfg Config) error {
	logger.Info("Initializing DataDir", "dir", cfg.DataDir)
	if err := utils.CreatePathIfNotExists(cfg.DataDir); err != nil {
		return err
	}
	
	if err := consensus.InitNode(cfg.consensusCfg); err != nil {
		return err
	}

	
	if err := database.Init(cfg.dbCfg); err != nil {
		return err
	}

	return nil
}
