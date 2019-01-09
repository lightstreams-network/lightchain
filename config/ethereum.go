package config

import (
	"path"
	"gopkg.in/urfave/cli.v1"
	"os"
	"time"
	"path/filepath"
	"io/ioutil"
	
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	ethCore "github.com/ethereum/go-ethereum/core"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	
	"github.com/lightstreams-network/lightchain/utils"
)

// General settings
var GenesisPathFlag = ethUtils.DirectoryFlag{
	Name:  "genesis",
	Usage: "Genesis path",
}


type ethstatsConfig struct {
	URL string `toml:",omitempty"`
}

type GethConfig struct {
	Eth      eth.Config
	Node     node.Config
	Ethstats ethstatsConfig
}


// DefaultNodeConfig returns the default configuration for a go-ethereum node
func DefaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.Version
	cfg.HTTPModules = append(cfg.HTTPModules, "eth")
	cfg.WSModules = append(cfg.WSModules, "eth")
	cfg.IPCPath = "geth.ipc"
	cfg.DataDir = utils.DefaultHomeDir()

	return cfg
}

// SetNodeConfig takes a node configuration and applies lightchain specific configuration
func SetNodeConfig(ctx *cli.Context, cfg *node.Config) {
	cfg.P2P.MaxPeers = 0
	cfg.P2P.NoDiscovery = true
	if ctx.GlobalIsSet(utils.HomeDirFlag.Name) {
		cfg.DataDir = ctx.GlobalString(utils.HomeDirFlag.Name)
	}
}

// SetEthConfig takes a ethereum configuration and applies lightchain specific configuration
func SetEthConfig(cfg *eth.Config) {
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
	homeDir := MakeHomeDir(ctx)
	genesisPath := path.Join(homeDir, "genesis.json")

	if ctx.GlobalIsSet(GenesisPathFlag.Name) {
		genesisPath = ctx.GlobalString(GenesisPathFlag.Name)
	}

	return genesisPath
}

func ReadGenesisFile(ctx *cli.Context) (*ethCore.Genesis, error) {
	genesisPath := MakeGenesisPath(ctx)
	ethLog.Info("Trying to reading genesis", "dir", genesisPath)
	genesisBlob, err := utils.ReadFileContent(genesisPath)
	if err != nil {
		ethLog.Warn("Cannot read genesis. Using default", err)
		genesisBlob, err = readDefaultGenesis()
		if err != nil {
			return nil, err
		}
	}
	
	genesis, err := utils.ParseBlobGenesis(genesisBlob)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}


func readDefaultGenesis() ([]byte, error) {
	fPath, err := filepath.Abs(filepath.Join(projectRootPath, "setup/genesis.json"))
	if err != nil {
		return nil, err
	}
	return utils.ReadFileContent(fPath)
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
