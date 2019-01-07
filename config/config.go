package config

import (
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
	"fmt"
	"encoding/json"
)

const (
	// Client identifier to advertise over the network
	clientIdentifier = "lightchain"
	// Environment variable for home dir
	emHome = "EMHOME"
)

var (
	// GenesisTargetGasLimit is the target gas limit of the Genesis block.
	// #unstable
	GenesisTargetGasLimit = uint64(100000000)
)

// General settings
var GenesisPathFlag = ethUtils.DirectoryFlag{
	Name:  "genesis",
	Usage: "Genesis path",
}

// The config folder name inside the lightchain's node --datadir ~/.lightchain
const ConfigFolderName = "lightchain"

// The config filename inside the ConfigFolderName
const ConfigFilename = "lightstreams_config.json"

type Config struct {
	DeploymentWhitelist []string `json:"deploymentWhitelist"`
}


var projectRootPath = filepath.Join(os.Getenv("GOPATH"), "src/github.com/lightstreams-network", "lightchain")

type ethstatsConfig struct {
	URL string `toml:",omitempty"`
}

type GethConfig struct {
	Eth      eth.Config
	Node     node.Config
	Ethstats ethstatsConfig
}

func NewLightchainConfig(path string) (Config, error) {
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("unable to open LS's configuration file %s", path)
	}

	cfg := new(Config)
	err = json.Unmarshal(buffer, cfg)
	if err != nil {
		return Config{}, err
	}

	return *cfg, nil
}

func ReadDefaultLsConfigBlob() ([]byte) {
	return []byte(`
{
    "deploymentWhitelist": [""]
}`)
}


// makeDataDir retrieves the currently requested data directory
func MakeDataDir(ctx *cli.Context) string {
	dPath := node.DefaultDataDir()

	emHome := os.Getenv(emHome)
	if emHome != "" {
		dPath = emHome
	}

	if ctx.GlobalIsSet(ethUtils.DataDirFlag.Name) {
		dPath = ctx.GlobalString(ethUtils.DataDirFlag.Name)
	}

	if dPath == "" {
		ethUtils.Fatalf("Cannot determine default data directory, please set manually (--datadir)")
	}

	return dPath
}

// DefaultNodeConfig returns the default configuration for a go-ethereum node
// #unstable
func DefaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.Version
	cfg.HTTPModules = append(cfg.HTTPModules, "eth")
	cfg.WSModules = append(cfg.WSModules, "eth")
	cfg.IPCPath = "geth.ipc"

	emHome := os.Getenv(emHome)
	if emHome != "" {
		cfg.DataDir = emHome
	}

	return cfg
}

// SetLightchainNodeDefaultConfig takes a node configuration and applies lightchain specific configuration
func SetLightchainNodeDefaultConfig(cfg *node.Config) {
	cfg.P2P.MaxPeers = 0
	cfg.P2P.NoDiscovery = true
}

// SetLightchainEthDefaultConfig takes a ethereum configuration and applies lightchain specific configuration
func SetLightchainEthDefaultConfig(cfg *eth.Config) {
	// PoW is being replaced by PoA with the usage of Tendermint
	cfg.Ethash.PowMode = ethash.ModeFake
	
	// Due to the low usages of the blockchain we need to reduce cache size to prevent huge number
	// block replies on every restart. 
	// @TODO once the usage of blockchain is larger we should tune these values again accordingly
	cfg.DatabaseCache  = 32; // MB
	cfg.TrieCleanCache = 8;  // MB
	cfg.TrieDirtyCache = 0;  // MB
	cfg.TrieTimeout = 5 * time.Minute;
}

func MakeGenesisPath(ctx *cli.Context) string {
	genesisPath := ctx.Args().First()
	if genesisPath != "" {
		return genesisPath
	} else if ctx.GlobalIsSet(GenesisPathFlag.Name) {
		genesisPath = ctx.GlobalString(GenesisPathFlag.Name)
	} else {
		lightchainDataDir := MakeDataDir(ctx)
		genesisPath = path.Join(lightchainDataDir, "genesis.json")
	}

	return genesisPath
}

func ReadGenesisPath(genesisPath string) ([]byte, error) {
	genesisBlob, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		return nil, err
	}

	return genesisBlob, nil
}

func ReadDefaultGenesis() ([]byte, error) {
	fPath, err := filepath.Abs(filepath.Join(projectRootPath, "setup/genesis.json"))
	if err != nil {
		return nil, err
	}
	return ReadGenesisPath(fPath)
}

func ReadDefaultKeystore() (map[string][]byte, error) {
	dPath, err := filepath.Abs(filepath.Join(projectRootPath, "setup/keystore"))
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

func ReadTendermintDefaultGenesis() ([]byte, error) {
	fPath, err := filepath.Abs(filepath.Join(projectRootPath, "setup/tendermint/genesis.json"))
	if err != nil {
		return nil, err
	}
	return ReadGenesisPath(fPath)
}

func ReadTendermintDefaultConfig() ([]byte, error) {
	fPath, err := filepath.Abs(filepath.Join(projectRootPath, "setup/tendermint/config.toml"))
	if err != nil {
		return nil, err
	}
	return ReadGenesisPath(fPath)
}
