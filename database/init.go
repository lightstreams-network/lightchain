package database

import (
	"path/filepath"
	"os"
	"encoding/json"
	"reflect"
	"fmt"
	"io/ioutil"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/lightstreams-network/lightchain/setup"
	ethCore "github.com/ethereum/go-ethereum/core"
	tmtLog "github.com/tendermint/tendermint/libs/log"
	"github.com/lightstreams-network/lightchain/log"
)

func Init(cfg Config, ntw setup.Network) error {
	logger := log.NewLogger().With("engine", "database")
	keystoreDir := cfg.keystoreDir()
	if err := os.MkdirAll(keystoreDir, os.ModePerm); err != nil {
		return err
	}

	var keystoreFiles map[string][]byte
	var genesisBlob []byte
	var err error
	switch ntw {
	case setup.SiriusNetwork:
		if keystoreFiles, err = setup.ReadSiriusDatabaseKeystore(); err != nil {
			return err
		}
		if genesisBlob, err = setup.ReadSiriusDatabaseGenesis(); err != nil {
			return err
		}
	case setup.StandaloneNetwork:
		if keystoreFiles, err = setup.ReadStandaloneDatabaseKeystore(); err != nil {
			return err
		}
		if genesisBlob, err = setup.ReadStandaloneDatabaseGenesis(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid network selected %s", ntw)
	}
	
	if err = writeKeystoreFiles(logger, keystoreDir, keystoreFiles); err != nil {
		err = fmt.Errorf("could not open write keystore: %v", err)
		return err
	}
	
	genesis, err := parseBlobGenesis(genesisBlob)
	if err != nil {
		err = fmt.Errorf("reading genesis err: %v", err)
		return err
	}

	if err = writeGenesisFile(cfg.genesisPath(), genesis); err != nil {
		err = fmt.Errorf("could not write genesis file: %v", err)
		return err
	}
	logger.Info("Generated genesis block", "path", cfg.genesisPath())

	chainDataDir := cfg.chainDbDir()
	chainDb, err := ethdb.NewLDBDatabase(chainDataDir, 0, 0)
	_, hash, err := ethCore.SetupGenesisBlock(chainDb, genesis)
	if err != nil {
		err = fmt.Errorf("failed to write genesis block: %v", err)
		return err
	}
	logger.Info("Successfully wrote genesis block and/or chain rule set", "hash", hash)
	
	return nil
}

func readGenesisFile(genesisPath string) (*ethCore.Genesis, error) {
	genesisBlob, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		return nil, err
	}

	genesis, err := parseBlobGenesis(genesisBlob)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}

func writeGenesisFile(genesisPath string, genesis *ethCore.Genesis) error {
	genesisBlob, err := genesis.MarshalJSON()
	if err != nil {
		return err
	}
	
	f, err := os.Create(genesisPath)
	if err != nil {
		return err
	}

	if _, err := f.Write(genesisBlob); err != nil {
		return err
	}
	
	return nil
}

func writeKeystoreFiles(logger tmtLog.Logger, keystoreDir string, keystoreFiles map[string][]byte) error {
	for filename, content := range keystoreFiles {
		storeFileName := filepath.Join(keystoreDir, filename)
		f, err := os.Create(storeFileName)
		if err != nil {
			logger.Error("Cannot create file", storeFileName, err)
			continue
		}
		if _, err := f.Write(content); err != nil {
			logger.Error("write content %q err: %v", storeFileName, err)
		}
		if err := f.Close(); err != nil {
			return err
		}

		logger.Info("Successfully wrote keystore files", "keystore", storeFileName)
	}

	return nil
}

// parseGenesisOrDefault tries to read the content from provided
// genesisPath. If the path is empty or doesn't exist, it will
// use defaultGenesisBytes as the fallback genesis source. Otherwise,
// it will open that path and if it encounters an error that doesn't
// satisfy os.IsNotExist, it returns that error.
func parseBlobGenesis(genesisBlob []byte) (*ethCore.Genesis, error) {
	genesis := new(ethCore.Genesis)
	if err := json.Unmarshal(genesisBlob, genesis); err != nil {
		return nil, err
	}

	if reflect.DeepEqual(blankGenesis, genesis) {
		return nil, errBlankGenesis
	}

	return genesis, nil
}
