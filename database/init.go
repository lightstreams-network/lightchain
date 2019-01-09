package database

import (
	"path/filepath"
	"os"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"gopkg.in/urfave/cli.v1"
	
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/lightstreams-network/lightchain/utils"
)

func Init(ctx *cli.Context) error {
	// Step 1: Init chain within --datadir by read genesis
	chainDataDir := utils.MakeChainDataDir(ctx)
	chainDb, err := ethdb.NewLDBDatabase(chainDataDir, 0, 0)
	if err != nil {
		ethUtils.Fatalf("could not open database: %v", err)
	}

	keystoreDir := utils.MakeKeystoreDir(ctx)
	if err := os.MkdirAll(keystoreDir, os.ModePerm); err != nil {
		ethUtils.Fatalf("mkdirAll keyStoreDir: %v", err)
	}

	keystoreCfg, err := ReadDefaultKeystore()
	if err != nil {
		ethUtils.Fatalf("could not open read keystore: %v", err)
	}

	for filename, content := range keystoreCfg {
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

	genesis, err := ReadGenesisFile(ctx)
	if err != nil {
		ethLog.Warn("reading genesis err: %v", err)
	}

	genesisBlob, err := genesis.MarshalJSON()
	if err != nil {
		ethLog.Warn("marshaling content err: %v", err)
	}
	
	genesisFileName := MakeGenesisPath(ctx)
	f, err := os.Create(genesisFileName)
	if _, err := f.Write(genesisBlob); err != nil {
		ethLog.Warn("write content %q err: %v", genesisFileName, err)
	} else {
		ethLog.Info("Using genesis block", "block", genesis)
	}
	
	_, hash, err := core.SetupGenesisBlock(chainDb, genesis)
	if err != nil {
		ethUtils.Fatalf("failed to write genesis block: %v", err)
		return err
	}

	ethLog.Info("Successfully wrote genesis block and/or chain rule set", "hash", hash)
	return nil
}
