package consensus

import (
	"path/filepath"
	"github.com/spf13/viper"
	"gopkg.in/urfave/cli.v1"
	"path"
	"fmt"
	"os"
	
	tmtConfig "github.com/tendermint/tendermint/config"
	
	"github.com/lightstreams-network/lightchain/utils"
	"io/ioutil"
)


var (
	DataDirPath = "consensus"
)

type TendermintConfig struct{
	RpcListenPort uint16 	// Default: 26657
	P2pListenPort uint16 	// Default: 26656
	ProxyListenPort uint16 	// Default: 26658
}

type Config struct {
	DataDir         string
	RpcListenPort   uint16 // Default: 26657
	P2pListenPort   uint16 // Default: 26656
	ProxyListenPort uint16 // Default: 26658
	tmtCfg          *tmtConfig.Config
}

func NewConfig(dataDir string, rpcListenPort uint16, p2pListenPort uint16, proxyListenPort uint16) Config {
	defaultCfg := tmtConfig.DefaultConfig()
	defaultCfg.SetRoot(dataDir)
	return Config {
		dataDir,
		rpcListenPort,
		p2pListenPort,
		proxyListenPort,
		defaultCfg,
	}
}

func (c Config) TendermintConfigPath () string {
	return filepath.Join(c.DataDir, "config")
}

func (c Config) TendermintConfigFilePath () string {
	return filepath.Join(c.TendermintConfigPath(), "config.toml")
}

func ensureTendermintDataDir(c Config) error {
	tmtConfig.EnsureRoot(c.DataDir)
	return nil
}


func MakeTendermintDir(ctx *cli.Context) string {
	homeDir := utils.MakeHomeDir(ctx)
	dataDir := filepath.Join(homeDir, "tmtCfg")
	tmtConfig.EnsureRoot(dataDir)
	return dataDir
}

func MakeTendermintConfig(ctx *cli.Context) TendermintConfig {
	return TendermintConfig {
		uint16(ctx.GlobalInt(utils.TendermintRpcListenPortFlag.Name)),
		uint16(ctx.GlobalInt(utils.TendermintP2PListenPortFlag.Name)),
		uint16(ctx.GlobalInt(utils.TendermintProxyListenPortFlag.Name)),
	}
}

func ParseTendermintConfig(ctx *cli.Context) (*tmtConfig.Config, error) {
	cfg := tmtConfig.DefaultConfig()
	dataDir := MakeTendermintDir(ctx);
	cfg.SetRoot(dataDir)
	initViper(path.Join(dataDir, "config"))
	err := viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	
	if err = cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("Error in config file: %v", err)
	}
	
	return cfg, err
}


//@TODO Migrate into constant
func readTendermintDefaultGenesis() ([]byte, error) {
	fPath, err := filepath.Abs(filepath.Join(utils.ProjectRootPath, "setup/tendermint/genesis.json"))
	if err != nil {
		return nil, err
	}
	
	return ioutil.ReadFile(fPath)
}

//@TODO Migrate into constant
func readTendermintDefaultConfig() ([]byte, error) {
	fPath, err := filepath.Abs(filepath.Join(utils.ProjectRootPath, "setup/tendermint/config.toml"))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(fPath)
}

func initViper(cfgPath string) error {
	viper.AddConfigPath(cfgPath) // search root directory /config
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// stderr, so if we redirect output to json file, this doesn't appear
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		// ignore not found error, return other errors
		return err
	}

	return nil
}