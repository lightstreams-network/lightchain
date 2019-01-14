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
	ctx	*cli.Context
}

type ethstatsConfig struct {
	URL string `toml:",omitempty"`
}

type GethConfig struct {
	Eth      eth.Config
	Node     ethNode.Config
	Ethstats ethstatsConfig
}


func NewConfig(dataDir string, ctx *cli.Context) (Config, error) {
	gethCfg := GethConfig{
		Eth:  eth.DefaultConfig,
		Node: DefaultEthNodeConfig(dataDir),
	}

	ethUtils.SetNodeConfig(ctx, &gethCfg.Node)
	setNodeDefaultConfig(&gethCfg.Node, dataDir)
	
	cfg := Config{
		dataDir,
		gethCfg,
		ctx,
	}

	//ethereumNode, err := ethNode.New(&cfg.GethConfig.Node)
	//if err != nil {
	//	return cfg, err
	//}
	//
	//ethUtils.SetEthConfig(ctx, ethereumNode, &cfg.GethConfig.Eth)
	//setEthDefaultConfig(&cfg.GethConfig.Eth)
	
	return cfg, nil
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

// SetNodeConfig takes a ethereum configuration and applies lightchain specific configuration
func setNodeDefaultConfig(cfg *ethNode.Config, dataDir string) {
	cfg.P2P.MaxPeers = 0
	cfg.P2P.NoDiscovery = true
	cfg.DataDir = dataDir
}

// SetEthConfig takes a ethereum configuration and applies lightchain specific configuration
func setEthDefaultConfig(cfg *eth.Config) {
	// PoW is being replaced by PoA with the usage of Tendermint
	cfg.Ethash.PowMode = ethash.ModeFake
	
	// Due to the low usages of the blockchain we need to reduce cache size to prevent huge number
	// block replies on every restart. 
	// @TODO once the usage of blockchain is larger we should tune these values again accordingly
	cfg.DatabaseCache  = 32 // MB
	cfg.TrieCleanCache = 8  // MB
	cfg.TrieDirtyCache = 0  // MB
	cfg.TrieTimeout = 5 * time.Minute
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
