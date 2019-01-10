package consensus

import (
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/p2p"
	
	tmtCommon "github.com/tendermint/tendermint/libs/common"
	"github.com/lightstreams-network/lightchain/log"
)

func InitNode(cfg Config) error {
	logger := log.NewLogger()
	if err := ensureTendermintDataDir(cfg); err != nil {
		return nil
	}

	privValFile := cfg.tmtCfg.PrivValidatorFile()
	var pv *privval.FilePV
	if tmtCommon.FileExists(privValFile) {
		pv = privval.LoadFilePV(privValFile)
		logger.Info("Found private validator", "path", privValFile)
	} else {
		pv = privval.GenFilePV(privValFile)
		pv.Save()
		logger.Info("Generated private validator", "path", privValFile)
	}

	nodeKeyFile := cfg.tmtCfg.NodeKeyFile()
	if tmtCommon.FileExists(nodeKeyFile) {
		logger.Info("Found node key", "path", nodeKeyFile)
	} else {
		if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
			return err
		}
		logger.Info("Generated node key", "path", nodeKeyFile)
	}

	// genesis file
	genFile := cfg.tmtCfg.GenesisFile()
	if tmtCommon.FileExists(genFile) {
		logger.Info("Found genesis file", "path", genFile)
	} else {
		genContent, err := readTendermintDefaultGenesis()
		if err != nil {
			return err
		}
		if err := tmtCommon.WriteFile(genFile, genContent, 0644); err != nil {
			return err
		}
		logger.Info("Generated genesis file", "path", genFile)
	}
	
	// Config file
	cfgFile := cfg.TendermintConfigPath()
	cfgDoc, err := readTendermintDefaultConfig()
	if err != nil {
		return err
	}
	if err := tmtCommon.WriteFile(cfgFile, cfgDoc, 0644); err != nil {
		return err
	}
	logger.Info("Generated config file", "path", cfgFile)
	
	return nil
}

