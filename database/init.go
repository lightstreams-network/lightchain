package database

import (
	"path/filepath"
	"os"
	ethCore "github.com/ethereum/go-ethereum/core"
	
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/ethdb"
)

func Init(cfg Config) error {
	keystoreDir := cfg.KeystoreDir()
	if err := os.MkdirAll(keystoreDir, os.ModePerm); err != nil {
		ethUtils.Fatalf("mkdirAll keyStoreDir: %v", err)
	}

	keystoreFiles, err := readDefaultKeystore()
	if err != nil {
		ethUtils.Fatalf("could not open read keystore: %v", err)
		return err
	}
	if err = writeKeystoreFiles(keystoreDir, keystoreFiles); err != nil {
		ethUtils.Fatalf("could not open write keystore: %v", err)
		return err
	}
	
	genesisPath := cfg.GenesisPath()
	genesis, err := readGenesisFile(genesisPath)
	if err != nil {
		ethLog.Warn("reading genesis err: %v", err)
	}
	if err = writeGenesisFile(genesisPath, genesis); err != nil {
		ethUtils.Fatalf("could not write genesis file: %v", err)
		return err
	}

	chainDataDir := cfg.ChainDbDir()
	chainDb, err := ethdb.NewLDBDatabase(chainDataDir, 0, 0)
	_, hash, err := ethCore.SetupGenesisBlock(chainDb, genesis)
	if err != nil {
		ethUtils.Fatalf("failed to write genesis block: %v", err)
		return err
	}

	ethLog.Info("Successfully wrote genesis block and/or chain rule set", "hash", hash)
	return nil
}

func writeGenesisFile(genesisPath string, genesis *ethCore.Genesis) error {
	genesisBlob, err := genesis.MarshalJSON()
	if err != nil {
		ethLog.Warn("marshaling content err: %v", err)
	}
	
	f, err := os.Create(genesisPath)
	if _, err := f.Write(genesisBlob); err != nil {
		ethLog.Warn("write content %q err: %v", genesisPath, err)
	} else {
		ethLog.Info("Using genesis block", "block", genesis)
	}
	
	return nil
}

func writeKeystoreFiles(keystoreDir string, keystoreFiles map[string][]byte) error {
	for filename, content := range keystoreFiles {
		storeFileName := filepath.Join(keystoreDir, filename)
		f, err := os.Create(storeFileName)
		if err != nil {
			ethLog.Error("Cannot create file", storeFileName, err)
			continue
		}
		if _, err := f.Write(content); err != nil {
			ethLog.Error("write content %q err: %v", storeFileName, err)
		}
		if err := f.Close(); err != nil {
			return err
		}

		ethLog.Info("Successfully wrote keystore files", "keystore", storeFileName)
	}

	return nil
}
