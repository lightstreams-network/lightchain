package consensus

import (
	"fmt"
	
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/p2p"
	tmtCommon "github.com/tendermint/tendermint/libs/common"
	tmtConfig "github.com/tendermint/tendermint/config"
	
	"github.com/lightstreams-network/lightchain/log"
	"github.com/lightstreams-network/lightchain/setup"
)


func Init(cfg Config, ntw setup.Network, logger log.Logger) error {
	// This is a necessary evil because
	// Tendermint is using panics instead of errors where they shouldn't...
	defer recoverNodeInitPanic()

	createConsensusDataDirIfNotExists(cfg.dataDir)

	privateValidatorFile := cfg.tendermintCfg.PrivValidatorFile()
	var pv *privval.FilePV
	if tmtCommon.FileExists(privateValidatorFile) {
		pv = privval.LoadFilePV(privateValidatorFile)
		logger.Info("Found private validator", "path", privateValidatorFile)
	} else {
		pv = privval.GenFilePV(privateValidatorFile)
		pv.Save()
		logger.Info("Generated private validator", "path", privateValidatorFile)
	}

	nodeKeyFile := cfg.tendermintCfg.NodeKeyFile()
	if tmtCommon.FileExists(nodeKeyFile) {
		logger.Info("Found node key", "path", nodeKeyFile)
	} else {
		if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
			return err
		}
		logger.Info("Generated node key", "path", nodeKeyFile)
	}

	var genContent []byte
	var cfgDoc []byte
	var err error
	cfgFilePath := cfg.TendermintConfigFilePath()
	genFile := cfg.tendermintCfg.GenesisFile()

	switch ntw {
	case setup.SiriusNetwork:
		if genContent, err = setup.ReadSiriusConsensusGenesis(); err != nil {
			return err
		}
		if cfgDoc, err = setup.ReadSiriusConsensusConfig(); err != nil {
			return err
		}
	case setup.StandaloneNetwork:
		if genContent, err = setup.CreateStandaloneConsensusGenesis(pv); err != nil {
			return err
		}
		if cfgDoc, err = setup.ReadStandaloneConsensusConfig(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Invalid network selected %s", ntw)
	}
	
	if tmtCommon.FileExists(genFile) {
		logger.Info("Found genesis file", "path", genFile)
	} else {
		if err := tmtCommon.WriteFile(genFile, genContent, 0644); err != nil {
			return err
		}
	}
	
	if err := tmtCommon.WriteFile(cfgFilePath, cfgDoc, 0644); err != nil {
		return err
	}
	
	return nil
}

func createConsensusDataDirIfNotExists(dataDir string) {
	tmtConfig.EnsureRoot(dataDir)
}

func recoverNodeInitPanic() error {
	if r := recover(); r != nil {
		return fmt.Errorf("panic occured initializing consensus node init. %v", r)
	}

	return nil
}
