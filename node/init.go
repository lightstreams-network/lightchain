package node

import (
	"os"
	
	"github.com/lightstreams-network/lightchain/log"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/setup"
)

func InitNode(cfg Config, ntw setup.Network) error {
	var logger = log.NewLogger()
	logger.With("module", "node")
	logger.Info("Initializing lightchain node data dir...", "dir", cfg.DataDir)

	if err := os.MkdirAll(cfg.DataDir, os.ModePerm); err != nil {
		return err
	}
	
	if err := consensus.Init(cfg.consensusCfg, ntw, logger); err != nil {
		return err
	}

	if err := database.Init(cfg.dbCfg, ntw, logger); err != nil {
		return err
	}

	return nil
}
