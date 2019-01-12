package database

import (
	"os"
	"time"
	"path/filepath"
	"io/ioutil"
	"errors"
	"reflect"
	"gopkg.in/urfave/cli.v1"

	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	ethCore "github.com/ethereum/go-ethereum/core"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/eth"
	ethNode "github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	
	"github.com/lightstreams-network/lightchain/utils"
	"encoding/json"
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
	Eth      eth.Config
	Node     ethNode.Config
	Ethstats ethstatsConfig
}

func NewConfig(dataDir string, ethCfg eth.Config, nodeCfg ethNode.Config, ethUrl string) (Config) {
	gethCfg := GethConfig {
		ethCfg, nodeCfg, ethstatsConfig { ethUrl },
	}

	return Config{
		dataDir,
		gethCfg,
	}
}

func NewConfigNode(dataDir string, ctx *cli.Context) (Config, error) {
	gethCfg := GethConfig{
		Eth:  eth.DefaultConfig,
		Node: DefaultEthNodeConfig(),
	}

	ethUtils.SetNodeConfig(ctx, &gethCfg.Node)
	setNodeDefaultConfig(&gethCfg.Node, dataDir)
	
	cfg := Config{
		dataDir,
		gethCfg,
	}

	ethereumNode, err := ethNode.New(&cfg.GethConfig.Node)
	if err != nil {
		return cfg, err
	}

	ethUtils.SetEthConfig(ctx, ethereumNode, &cfg.GethConfig.Eth)
	setEthDefaultConfig(&cfg.GethConfig.Eth)

 	// @TODO Review the need of including `stack` as part of the method output
	return cfg, nil
}



// DefaultEthNodeConfig returns the default configuration for a go-ethereum ethereum
func DefaultEthNodeConfig() ethNode.Config {
	cfg := ethNode.DefaultConfig
	cfg.Name = utils.ClientIdentifier
	cfg.Version = params.Version
	cfg.HTTPModules = append(cfg.HTTPModules, "eth")
	cfg.WSModules = append(cfg.WSModules, "eth")
	cfg.IPCPath = "geth.ipc"
	cfg.DataDir = utils.DefaultDataDir()

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

func readGenesisFile(genesisPath string) (*ethCore.Genesis, error) {
	ethLog.Info("Trying to reading genesis", "dir", genesisPath)
	genesisBlob, err := utils.ReadFileContent(genesisPath)
	if err != nil {
		ethLog.Warn("Cannot read genesis. Using default", err)
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
	fPath, err := filepath.Abs(filepath.Join(utils.ProjectRootPath, "setup/genesis.json"))
	if err != nil {
		return nil, err
	}
	return utils.ReadFileContent(fPath)
}

func readDefaultKeystore() (map[string][]byte, error) {
	dPath, err := filepath.Abs(filepath.Join(utils.ProjectRootPath, "setup/keystore"))
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
