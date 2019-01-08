package tendermint

import (
	"os"
	"path/filepath"
	"gopkg.in/urfave/cli.v1"
	
	tmtCfg "github.com/tendermint/tendermint/config"
	tmtCommon "github.com/tendermint/tendermint/libs/common"
	tmtNode "github.com/tendermint/tendermint/node"
	tmtLog "github.com/tendermint/tendermint/libs/log"
	
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	
	"github.com/lightstreams-network/lightchain/config"
	"github.com/lightstreams-network/lightchain/utils"
	"fmt"
	"os/signal"
	"syscall"
	"github.com/tendermint/tendermint/proxy"
)

func InitNode(ctx *cli.Context) error {
	cfg := tmtCfg.DefaultConfig()
	logger := utils.LightchainLogger()
	
	// Step 1: Init chain within --datadir by read genesis
	dataDir := config.MakeTendermintDir(ctx)
	cfg.SetRoot(dataDir)

	privValFile := cfg.PrivValidatorFile()
	var pv *privval.FilePV
	if tmtCommon.FileExists(privValFile) {
		pv = privval.LoadFilePV(privValFile)
		logger.Info("Found private validator", "path", privValFile)
	} else {
		pv = privval.GenFilePV(privValFile)
		pv.Save()
		logger.Info("Generated private validator", "path", privValFile)
	}

	nodeKeyFile := cfg.NodeKeyFile()
	if tmtCommon.FileExists(nodeKeyFile) {
		logger.Info("Found node key", "path", nodeKeyFile)
	} else {
		if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
			return err
		}
		logger.Info("Generated node key", "path", nodeKeyFile)
	}

	// genesis file
	genFile := cfg.GenesisFile()
	if tmtCommon.FileExists(genFile) {
		logger.Info("Found genesis file", "path", genFile)
	} else {
		genDoc, err := config.ReadTendermintDefaultGenesis()
		if err != nil {
			return err
		}
		if err := tmtCommon.WriteFile(genFile, genDoc, 0644); err != nil {
			return err
		}
		logger.Info("Generated genesis file", "path", genFile)
	}
	
	// Config file
	cfgDir := filepath.Join(dataDir, "config")
	cfgFile := filepath.Join(cfgDir, "config.toml")
	cfgDoc, err := config.ReadTendermintDefaultConfig()
	if err != nil {
		return err
	}
	if err := tmtCommon.WriteFile(cfgFile, cfgDoc, 0644); err != nil {
		return err
	}
	logger.Info("Generated config file", "path", cfgFile)
	
	return nil
}


func StartNode(ctx *cli.Context, n *tmtNode.Node) error {
	logger := utils.LightchainLogger()
	
	// Stop upon receiving SIGTERM or CTRL-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			logger.Error(fmt.Sprintf("captured %v, exiting...", sig))
			if n.IsRunning() {
				n.Stop()
			}
			os.Exit(1)
		}
	}()

	if err := n.Start(); err != nil {
		return fmt.Errorf("Failed to start node: %v", err)
	}
	logger.Info("Started node", "nodeInfo", n.Switch().NodeInfo())

	// Run forever
	select {}

	return nil
}



func CreateNewNode(ctx *cli.Context, logger tmtLog.Logger) (*tmtNode.Node, error) {
	tmtConfig := config.MakeTendermintConfig(ctx)
	cfg, err := config.ParseTendermintConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Generate node PrivKey
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}
	
	defaultRPCfg := tmtCfg.DefaultRPCConfig()
	defaultRPCfg.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", tmtConfig.RpcListenPort)
	cfg.RPC = defaultRPCfg
	
	defaultP2P := tmtCfg.DefaultP2PConfig()
	defaultP2P.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", tmtConfig.P2pListenPort)
	cfg.P2P = defaultP2P
	
	return tmtNode.NewNode(cfg,
		privval.LoadOrGenFilePV(cfg.PrivValidatorFile()),
		nodeKey,
		proxy.DefaultClientCreator(cfg.ProxyApp, cfg.ABCI, cfg.DBDir()),
		tmtNode.DefaultGenesisDocProviderFunc(cfg),
		tmtNode.DefaultDBProvider,
		tmtNode.DefaultMetricsProvider(cfg.Instrumentation),
		logger,
	)
}

