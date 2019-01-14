package database

import (
	"path/filepath"
	"os"
	
	ethCore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	
	"github.com/lightstreams-network/lightchain/log"
)

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
