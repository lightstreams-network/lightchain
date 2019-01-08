package config

import (
	"path/filepath"
	"github.com/spf13/viper"
	"gopkg.in/urfave/cli.v1"
	"path"
	"fmt"
	"os"
	
	tmtCfg "github.com/tendermint/tendermint/config"
	
	"github.com/lightstreams-network/lightchain/utils"
)

type TendermintConfig struct{
	RpcListenPort uint16 	// Default: 26657
	P2pListenPort uint16 	// Default: 26656
	ProxyListenPort uint16 	// Default: 26658
}

func MakeTendermintDir(ctx *cli.Context) string {
	homeDir := MakeHomeDir(ctx)
	dataDir := filepath.Join(homeDir, "tendermint")
	tmtCfg.EnsureRoot(dataDir)
	return dataDir
}

func MakeTendermintConfig(ctx *cli.Context) TendermintConfig {
	return TendermintConfig {
		uint16(ctx.GlobalInt(utils.TendermintRpcListenPortFlag.Name)),
		uint16(ctx.GlobalInt(utils.TendermintP2PListenPortFlag.Name)),
		uint16(ctx.GlobalInt(utils.ProxyListenPortFlag.Name)),
	}
}

func ParseTendermintConfig(ctx *cli.Context) (*tmtCfg.Config, error) {
	cfg := tmtCfg.DefaultConfig()
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
func ReadTendermintDefaultGenesis() ([]byte, error) {
	fPath, err := filepath.Abs(filepath.Join(projectRootPath, "setup/tendermint/genesis.json"))
	if err != nil {
		return nil, err
	}
	return utils.ReadFileContent(fPath)
}

//@TODO Migrate into constant
func ReadTendermintDefaultConfig() ([]byte, error) {
	fPath, err := filepath.Abs(filepath.Join(projectRootPath, "setup/tendermint/config.toml"))
	if err != nil {
		return nil, err
	}
	return utils.ReadFileContent(fPath)
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
