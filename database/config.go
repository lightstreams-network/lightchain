package database

import (
	"time"
	"path/filepath"
	"errors"
	"gopkg.in/urfave/cli.v1"

	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	ethCore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	ethNode "github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	
)

var (
	DataDirPath = "database"
	KeystorePath = "keystore"
	ChainDbPath = "chaindb"
	GenesisPath = "genesis.json"
	blankGenesis = new(ethCore.Genesis)
	errBlankGenesis = errors.New("could not parse a valid/non-blank Genesis")
)

type Config struct {
	DataDir string
	GethConfig GethConfig
}

type ethstatsConfig struct {
	URL string `toml:",omitempty"`
}

type GethConfig struct {
	EthCfg   eth.Config
	NodeCfg  ethNode.Config
	Ethstats ethstatsConfig
}


func NewConfig(dataDir string, ctx *cli.Context) (Config, error) {
	gethCfg := GethConfig{
		EthCfg:  eth.DefaultConfig,
		NodeCfg: DefaultEthNodeConfig(dataDir),
	}

	// Configure ethereum node server
	ethUtils.SetNodeConfig(ctx, &gethCfg.NodeCfg)
	gethCfg.NodeCfg.P2P.MaxPeers = 0
	gethCfg.NodeCfg.P2P.NoDiscovery = true
	gethCfg.NodeCfg.DataDir = dataDir
	
	// REMINDER: Following initialization is required to complete the configuration of the ethereum db
	ethereum, err := ethNode.New(&gethCfg.NodeCfg)
	if err != nil {
		return Config{}, err
	}

	// Configure ethereum db settings
	ethUtils.SetEthConfig(ctx, ethereum, &gethCfg.EthCfg)
	gethCfg.EthCfg.Ethash.PowMode = ethash.ModeFake
	
	// Due to the low usages of the blockchain we need to reduce cache size to prevent huge number block replies on every restart. 
	// @TODO once the usage of blockchain is larger we should tune these values again accordingly
	gethCfg.EthCfg.DatabaseCache  = 32 // MB
	gethCfg.EthCfg.TrieCleanCache = 8  // MB
	gethCfg.EthCfg.TrieDirtyCache = 0  // MB
	gethCfg.EthCfg.TrieTimeout = 5 * time.Minute
	
	gethCfg.EthCfg.Genesis, err = readGenesisFile(dataDir)
	if err != nil {
		return Config{}, err
	}
	
	return Config {
		dataDir,
		gethCfg,
	}, nil
}

// DefaultEthNodeConfig returns the default configuration for a go-ethereum ethereum
func DefaultEthNodeConfig(dataDir string) ethNode.Config {
	cfg := ethNode.DefaultConfig
	cfg.Name = "lightchain"
	cfg.Version = params.Version
	cfg.HTTPModules = append(cfg.HTTPModules, "eth")
	cfg.WSModules = append(cfg.WSModules, "eth")
	cfg.IPCPath = "geth.ipc"
	cfg.DataDir = dataDir

	return cfg
}

func (c Config) keystoreDir() string {
	return filepath.Join(c.DataDir, KeystorePath)
}

func (c Config) chainDbDir() string {
	return filepath.Join(c.DataDir, ChainDbPath)
}

func (c Config) genesisPath() string {
	return filepath.Join(c.DataDir, GenesisPath)
}
