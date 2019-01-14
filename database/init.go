package database

import (
	"path/filepath"
	"os"
	
	ethCore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	
	"github.com/lightstreams-network/lightchain/log"
	"io/ioutil"
	"encoding/json"
	"reflect"
)

var LightchainSetupDir = filepath.Join(os.Getenv("GOPATH"), "src/github.com/lightstreams-network", "lightchain", "setup")

func Init(cfg Config, logger log.Logger) error {
	keystoreDir := cfg.keystoreDir()
	if err := os.MkdirAll(keystoreDir, os.ModePerm); err != nil {
		logger.Error("mkdirAll keyStoreDir: %v", err)
	}

	keystoreFiles, err := readDefaultKeystore()
	if err != nil {
		logger.Error("could not open read keystore: %v", err)
		return err
	}
	if err = writeKeystoreFiles(logger, keystoreDir, keystoreFiles); err != nil {
		logger.Error("could not open write keystore: %v", err)
		return err
	}
	
	genesisPath := cfg.genesisPath()
	genesis, err := readGenesisFile(genesisPath)
	if err != nil {
		logger.Error("reading genesis err: %v", err)
		return err
	}
	if err = writeGenesisFile(genesisPath, genesis); err != nil {
		logger.Error("could not write genesis file: %v", err)
		return err
	}
	logger.Info("Generated genesis block", "path", genesisPath)

	chainDataDir := cfg.chainDbDir()
	chainDb, err := ethdb.NewLDBDatabase(chainDataDir, 0, 0)
	_, hash, err := ethCore.SetupGenesisBlock(chainDb, genesis)
	if err != nil {
		logger.Error("failed to write genesis block: %v", err)
		return err
	}

	logger.Info("Successfully wrote genesis block and/or chain rule set", "hash", hash)
	return nil
}

func readGenesisFile(genesisPath string) (*ethCore.Genesis, error) {
	genesisBlob, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		genesisBlob, err = readDefaultGenesis()
		if err != nil {
			return nil, err
		}
	}

	genesis, err := parseBlobGenesis(genesisBlob)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}

func readDefaultGenesis() ([]byte, error) {
	fPath, err := filepath.Abs(filepath.Join(LightchainSetupDir, "genesis.json"))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(fPath)
}

func readDefaultKeystore() (map[string][]byte, error) {
	dPath, err := filepath.Abs(filepath.Join(LightchainSetupDir, "keystore"))
	if err != nil {
		return nil, err
	}

	var files = make(map[string][]byte)
	err = filepath.Walk(dPath, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		files[info.Name()] = content
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files, nil
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

func writeKeystoreFiles(logger log.Logger, keystoreDir string, keystoreFiles map[string][]byte) error {
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
