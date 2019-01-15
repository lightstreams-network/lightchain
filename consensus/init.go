package consensus

import (
	"fmt"
	"path/filepath"
	"io/ioutil"
	"os"
	
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/p2p"
	tmtCommon "github.com/tendermint/tendermint/libs/common"
	tmtConfig "github.com/tendermint/tendermint/config"
	tmtime "github.com/tendermint/tendermint/types/time"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/go-amino"
	
	"github.com/lightstreams-network/lightchain/log"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var cdc = amino.NewCodec()

func init() {
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoRoute, nil)
}

var LightchainSetupDir = filepath.Join(os.Getenv("GOPATH"), "src/github.com/lightstreams-network", "lightchain", "setup")


func Init(cfg Config, ntw Network, logger log.Logger) error {
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
		var genContent []byte
		var err error
		if (ntw == SiriusNetwork) {
			if genContent, err = readSiriusDefaultGenesis(); err != nil {
				return err
			}
		} else {
			if genContent, err = newDefaultGenesis(pv); err != nil {
				return err
			}
		}

		if err := tmtCommon.WriteFile(genFile, genContent, 0644); err != nil {
			return err
		}

		logger.Info("Generated genesis file", "path", genFile)
	}

	cfgFilePath := cfg.TendermintConfigFilePath()
	if (ntw == SiriusNetwork) {
		cfgDoc, err := readSiriusDefaultConfig()
		if err != nil {
			return err
		}
	
		if err := tmtCommon.WriteFile(cfgFilePath, cfgDoc, 0644); err != nil {
			return err
		}
		logger.Info("Generated Tendermint config.toml", "path", cfgFilePath)
	}

	return nil
}

func createConsensusDataDirIfNotExists(dataDir string) {
	tmtConfig.EnsureRoot(dataDir)
}

func newDefaultGenesis(pv *privval.FilePV) ([]byte, error) {
	genDoc := types.GenesisDoc{
		ChainID:         fmt.Sprintf("test-chain-%v", tmtCommon.RandStr(6)),
		GenesisTime:     tmtime.Now(),
		ConsensusParams: types.DefaultConsensusParams(),
	}
	genDoc.Validators = []types.GenesisValidator{{
		Address: pv.GetPubKey().Address(),
		PubKey:  pv.GetPubKey(),
		Power:   10,
	}}
	
	genDocBytes, err := cdc.MarshalJSONIndent(genDoc, "", "  ")
	if err != nil {
		return nil, err
	}

	return genDocBytes, nil
}

func readSiriusDefaultGenesis() ([]byte, error) {
	fPath, err := filepath.Abs(filepath.Join(LightchainSetupDir, "tendermint/genesis.json"))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(fPath)
}

func readSiriusDefaultConfig() ([]byte, error) {
	fPath, err := filepath.Abs(filepath.Join(LightchainSetupDir, "tendermint/config.toml"))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(fPath)
}

func recoverNodeInitPanic() error {
	if r := recover(); r != nil {
		return fmt.Errorf("panic occured initializing consensus node init. %v", r)
	}

	return nil
}
