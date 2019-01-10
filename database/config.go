package database

import (
	"path"
	"gopkg.in/urfave/cli.v1"
	"os"
	"time"
	"path/filepath"
	"io/ioutil"
	"errors"
	"reflect"

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


// DefaultEthNodeConfig returns the default configuration for a go-ethereum ethNode
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

// SetNodeConfig takes a ethNode configuration and applies lightchain specific configuration
func SetNodeDefaultConfig(cfg *ethNode.Config, dataDir string) {
	cfg.P2P.MaxPeers = 0
	cfg.P2P.NoDiscovery = true
	cfg.DataDir = dataDir
}

// SetEthConfig takes a ethereum configuration and applies lightchain specific configuration
func SetEthDefaultConfig(cfg *eth.Config) {
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

func MakeGenesisPath(ctx *cli.Context) string {
	homeDir := utils.MakeHomeDir(ctx)
	genesisPath := path.Join(homeDir, "genesis.json")

	if ctx.GlobalIsSet(utils.GenesisPathFlag.Name) {
		genesisPath = ctx.GlobalString(utils.GenesisPathFlag.Name)
	}

	return genesisPath
}


func (c Config) KeystoreDir() string {
	return filepath.Join(c.DataDir, KeystorePath)
}

func (c Config) ChainDbDir() string {
	return filepath.Join(c.DataDir, ChainDbPath)
}

func (c Config) GenesisPath() string {
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
