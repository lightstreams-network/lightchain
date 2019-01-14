package consensus

import (
	"fmt"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/p2p"
	"github.com/lightstreams-network/lightchain/log"

	tmtCommon "github.com/tendermint/tendermint/libs/common"
	tmtConfig "github.com/tendermint/tendermint/config"
	"path/filepath"
	"io/ioutil"
	"os"
)

var LightchainSetupDir = filepath.Join(os.Getenv("GOPATH"), "src/github.com/lightstreams-network", "lightchain", "setup")

func Init(cfg Config, logger log.Logger) error {
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

	genFile := cfg.tendermintCfg.GenesisFile()
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

	cfgFilePath := cfg.TendermintConfigFilePath()
	cfgDoc, err := readTendermintDefaultConfig()
	if err != nil {
		return err
	}

	if err := tmtCommon.WriteFile(cfgFilePath, cfgDoc, 0644); err != nil {
		return err
	}
	logger.Info("Generated Tendermint config file", "path", cfgFilePath)

	return nil
}

func createConsensusDataDirIfNotExists(dataDir string) {
	tmtConfig.EnsureRoot(dataDir)
}

func readTendermintDefaultGenesis() ([]byte, error) {
	fPath, err := filepath.Abs(filepath.Join(LightchainSetupDir, "tendermint/genesis.json"))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(fPath)
}

func readTendermintDefaultConfig() ([]byte, error) {
	fPath, err := filepath.Abs(filepath.Join(LightchainSetupDir, "tendermint/config.toml"))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(fPath)
}

func recoverNodeInitPanic() error {
	if r := recover(); r!= nil {
		return fmt.Errorf("panic occured initializing consensus node init. %v", r)
	}

	return nil
}