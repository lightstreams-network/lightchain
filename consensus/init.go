package consensus

import (
	"fmt"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/p2p"
	
	tmtCommon "github.com/tendermint/tendermint/libs/common"
	tmtConfig "github.com/tendermint/tendermint/config"
	tmtDb "github.com/tendermint/tendermint/libs/db"
	tmtState "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/version"
	
	"github.com/lightstreams-network/lightchain/setup"
	"github.com/lightstreams-network/lightchain/log"
	stdtracer "github.com/lightstreams-network/lightchain/tracer"
)

func Init(cfg Config, ntw setup.Network, trcCfg stdtracer.Config) error {
	logger := log.NewLogger().With("engine", "consensus")

	createConsensusDataDirIfNotExists(cfg.dataDir)

	privateValidatorKeyFile := cfg.tendermintCfg.PrivValidatorKeyFile()
	privateValidatorStateFile := cfg.tendermintCfg.PrivValidatorStateFile()
	var pv *privval.FilePV
	if tmtCommon.FileExists(privateValidatorKeyFile) {
		pv = privval.LoadFilePV(privateValidatorKeyFile, privateValidatorStateFile)
		logger.Info("Found private validator key file", "path", privateValidatorKeyFile)
	} else {
		pv = privval.GenFilePV(privateValidatorKeyFile, privateValidatorStateFile)
		pv.Save()
		logger.Info("Generated private validator key file", "path", privateValidatorKeyFile)
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
	var protocolBlockVersion version.Protocol
	var err error
	cfgFilePath := cfg.TendermintConfigFilePath()
	genFile := cfg.tendermintCfg.GenesisFile()

	switch ntw {
	case setup.MainNetNetwork:
		protocolBlockVersion = mainNetProtocolBlockVersion
		if genContent, err = setup.ReadMainNetConsensusGenesis(); err != nil {
			return err
		}
		if cfgDoc, err = setup.ReadMainNetConsensusConfig(); err != nil {
			return err
		}
	case setup.SiriusNetwork:
		protocolBlockVersion = siriusProtocolBlockVersion
		if genContent, err = setup.ReadSiriusConsensusGenesis(); err != nil {
			return err
		}
		if cfgDoc, err = setup.ReadSiriusConsensusConfig(); err != nil {
			return err
		}
	case setup.StandaloneNetwork:
		protocolBlockVersion = standaloneProtocolBlockVersion
		if genContent, err = setup.CreateStandaloneConsensusGenesis(pv); err != nil {
			return err
		}
		if cfgDoc, err = setup.ReadStandaloneConsensusConfig(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid network selected %s", ntw)
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
	
	logger.Info("Initializing consensus StateDB", "database", cfg.tendermintCfg.DBDir())
	stateDB := tmtDb.NewDB("state", tmtDb.DBBackendType(cfg.tendermintCfg.DBBackend), cfg.tendermintCfg.DBDir())
	state, err := tmtState.LoadStateFromDBOrGenesisFile(stateDB, cfg.tendermintCfg.GenesisFile())
	state.Version.Consensus.Block = protocolBlockVersion
	logger.Info("Saving and closing consensus StateDB", "database", cfg.tendermintCfg.DBDir())
	tmtState.SaveState(stateDB, state)
	stateDB.Close()
	
	trc, err := newTracer(trcCfg)
	if err != nil {
		return err
	}
	
	trc.assertPersistedInitStateDb(cfg.tendermintCfg, ntw)

	return nil
}

func createConsensusDataDirIfNotExists(dataDir string) {
	tmtConfig.EnsureRoot(dataDir)
}