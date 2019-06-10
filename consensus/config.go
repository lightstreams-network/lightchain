package consensus

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/config"
	
)

const DataDirName = "consensus"

// Config is general consensus node config.
//
// Default values:
//
// - rpcListenPort 26657
// - p2pListenPort 26656
// - proxyAppName lightchain
type Config struct {
	dataDir       string
	tendermintCfg *config.Config
	rpcListenPort uint
	p2pListenPort uint
	proxyAppName  string
	metrics       bool
}

func NewConfig(dataDir string, rpcListenPort uint, p2pListenPort uint, proxyAppName string, metrics bool) (Config, error) {
	tendermintCfg := config.DefaultConfig()

	err := applyTendermintConfig(filepath.Join(dataDir, "config"), tendermintCfg)
	if err != nil {
		return Config{}, err
	}

	tendermintCfg.SetRoot(dataDir)
	tendermintCfg.RPC.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", rpcListenPort)
	tendermintCfg.P2P.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", p2pListenPort)
	tendermintCfg.ProxyApp = proxyAppName

	return Config{
		dataDir,
		tendermintCfg,
		rpcListenPort,
		p2pListenPort,
		proxyAppName,
		metrics,
	}, nil
}

func (c Config) TendermintConfigFilePath() string {
	return filepath.Join(filepath.Join(c.dataDir, "config"), "config.toml")
}

func (c Config) RPCListenPort() uint {
	return c.rpcListenPort
}


func (c Config) TendermintCfg() config.Config {
	return *c.tendermintCfg;
}

func applyTendermintConfig(configPath string, tmtCfg *config.Config) error {
	viper.AddConfigPath(configPath) // search root directory /config
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// stderr, so if we redirect output to json file, this doesn't appear
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		// ignore not found error, return other errors
		return err
	}

	if err := viper.Unmarshal(tmtCfg); err != nil {
		return err
	}

	if err := tmtCfg.ValidateBasic(); err != nil {
		return err
	}

	return nil
}
