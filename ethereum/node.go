package ethereum

import (
	"github.com/ethereum/go-ethereum/node"
	"github.com/lightstreams-network/lightchain/config"
	"path/filepath"
	"os"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"gopkg.in/urfave/cli.v1"
	
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	ethLog "github.com/ethereum/go-ethereum/log"
)

// Node is the main object.
type Node struct {
	node.Node
}

func InitNode(ctx *cli.Context) error {
	// Step 1: Init chain within --datadir by read genesis
	chainDataDir := config.MakeChainDataDir(ctx)
	chainDb, err := ethdb.NewLDBDatabase(chainDataDir, 0, 0)
	if err != nil {
		ethUtils.Fatalf("could not open database: %v", err)
	}

	keystoreDir := config.MakeKeystoreDir(ctx)
	if err := os.MkdirAll(keystoreDir, os.ModePerm); err != nil {
		ethUtils.Fatalf("mkdirAll keyStoreDir: %v", err)
	}

	keystoreCfg, err := config.ReadDefaultKeystore()
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

	genesis, err := config.ReadGenesisFile(ctx)
	if err != nil {
		ethLog.Warn("reading genesis err: %v", err)
	}

	genesisBlob, err := genesis.MarshalJSON()
	if err != nil {
		ethLog.Warn("marshaling content err: %v", err)
	}
	
	genesisFileName := config.MakeGenesisPath(ctx)
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

// New creates a new node.
func New(conf *node.Config) (*Node, error) {
	stack, err := node.New(conf)
	if err != nil {
		return nil, err
	}

	return &Node{*stack}, nil // nolint: vet
}

// Start starts base node and stop p2p server
func (n *Node) Start() error {
	// start p2p server
	err := n.Node.Start()
	if err != nil {
		return err
	}

	// Stop it Eth.p2p server
	n.Node.Server().Stop()

	return nil
}


func (n *Node) Stop() error {
	err := n.Node.Stop()
	if err != nil {
		return err
	}

	return nil
}
