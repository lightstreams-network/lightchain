package consensus

import (
	"gopkg.in/urfave/cli.v1"
	"fmt"
	
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	
	tmtConfig "github.com/tendermint/tendermint/config"
	tmtNode "github.com/tendermint/tendermint/node"
	tmtLog "github.com/tendermint/tendermint/libs/log"
	tmtFlags "github.com/tendermint/tendermint/libs/cli/flags"
	rpcTypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcClient "github.com/tendermint/tendermint/rpc/lib/client"
	
	"github.com/lightstreams-network/lightchain/log"
	"github.com/tendermint/tendermint/proxy"
)


type Node struct {
	tmtNode *tmtNode.Node
	cfg *Config
} 

func NewNode(cfg Config) (*Node, error) {
	return &Node{ nil, &cfg }, nil
}

func (n *Node) Start() error {
	return nil
}

func (n *Node) Stop() error {
	return nil
}

func (n *Node) NewURIClient() *rpcClient.URIClient {
	client := rpcClient.NewURIClient(n.cfg) // tendermint RPC client
	rpcTypes.RegisterAmino(client.Codec())
	return client
}

/// NEXT TO REFACTOR

func StartNode(ctx *cli.Context, n *tmtNode.Node) error {
	logger := log.NewLogger()

	if err := n.Start(); err != nil {
		return fmt.Errorf("Failed to start node: %v", err)
	}
	logger.Info("Started node", "nodeInfo", n.Switch().NodeInfo())

	return nil
}

func CreateNewNode(ctx *cli.Context, logger tmtLog.Logger) (*tmtNode.Node, error) {
	ctxTmtCfg := MakeTendermintConfig(ctx)
	cfg, err := ParseTendermintConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Generate node PrivKey
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}
	
	cfg.RPC.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", ctxTmtCfg.RpcListenPort)
	cfg.P2P.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", ctxTmtCfg.P2pListenPort)
	cfg.ProxyApp = fmt.Sprintf("tcp://127.0.0.1:%d", ctxTmtCfg.ProxyListenPort)
	
	logger, err = tmtFlags.ParseLogLevel(cfg.LogLevel, logger, tmtConfig.DefaultLogLevel())
	if err != nil {
		return nil, err
	}
	
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

