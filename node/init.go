package node

import (
	"os"
	
	"github.com/lightstreams-network/lightchain/log"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/network"
)

func Init(cfg Config, ntw network.Network) error {
	logger := log.NewLogger().With("engine", "node")
	logger.Info("Initializing lightchain node data dir...", "dir", cfg.DataDir)
	
	if err := os.MkdirAll(cfg.DataDir, os.ModePerm); err != nil {
		return err
	}
	
	if err := consensus.Init(cfg.consensusCfg, ntw, cfg.tracerCfg); err != nil {
		return err
	}

	if err := database.Init(cfg.dbCfg, ntw, cfg.tracerCfg); err != nil {
		return err
	}

	return nil
}
