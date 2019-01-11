package consensus

import (
	"path/filepath"
	"github.com/tendermint/tendermint/config"
	"fmt"
)

const DataDirName = "consensus"

// Config is general consensus node config.
//
// Default values:
//
// - rpcListenPort 26657
// - p2pListenPort 26656
// - ProxyListenPort 26658
//
type Config struct {
	dataDir         string
	tendermintCfg   *config.Config
	rpcListenPort   uint
	p2pListenPort   uint
	proxyListenPort uint
	proxyProtocol   string
}

func NewConfig(dataDir string, rpcListenPort uint, p2pListenPort uint, proxyListenPort uint, proxyProtocol string) Config {
	tendermintCfg := config.DefaultConfig()
	tendermintCfg.SetRoot(dataDir)
	tendermintCfg.RPC.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", rpcListenPort)
	tendermintCfg.P2P.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", p2pListenPort)
	tendermintCfg.ProxyApp = fmt.Sprintf("tcp://127.0.0.1:%d", proxyListenPort)

	return Config {
		dataDir,
		tendermintCfg,
		rpcListenPort,
		p2pListenPort,
		proxyListenPort,
		proxyProtocol,
	}
}

func (c Config) TendermintConfigFilePath() string {
	return filepath.Join(filepath.Join(c.dataDir, "config"), "config.toml")
}