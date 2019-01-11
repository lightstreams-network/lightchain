package node

import (
	"github.com/lightstreams-network/lightchain/utils"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/database"
)

func InitNode(cfg Config) error {
	var logger = log.NewLogger()
	logger.With("module", "node")
	logger.Info("Initializing lightchain node data dir...", "dir", cfg.DataDir)

	if err := utils.CreatePathIfNotExists(cfg.DataDir); err != nil {
		return err
	}
	
	if err := consensus.Init(cfg.consensusCfg, logger); err != nil {
		return err
	}

	if err := database.Init(cfg.dbCfg); err != nil {
		return err
	}

	return nil
}