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
	tendermint *tmtNode.Node
	cfg        *Config
	logger log.Logger
} 

func NewNode(cfg *Config) (*Node, error) {
	consensusLogger := log.NewLogger()
	consensusLogger.With("module", "consensus")
		
	return &Node{
		nil,
		cfg,
		consensusLogger,
	}, nil
}

func (n *Node) Start() error {
	logger := log.NewLogger()

	if err := n.tendermint.Start(); err != nil {
		return err
	}

	logger.Info("Started node", "nodeInfo", n.tendermint.Switch().NodeInfo())

	return nil
}

func (n *Node) Stop() error {
	return nil
}

func (n *Node) NewURIClient() *rpcClient.URIClient {
	lAddr := fmt.Sprintf("tcp://%s:%d", "127.0.0.1", n.cfg.RpcListenPort)
	client := rpcClient.NewURIClient(lAddr)
	rpcTypes.RegisterAmino(client.Codec())
	return client
}

/// NEXT TO REFACTOR

func CreateNewNode(ctx *cli.Context, logger tmtLog.Logger) (*tmtNode.Node, error) {

	// Generate node PrivKey
	nodeKey, err := p2p.LoadOrGenNodeKey(n.cfg.NodeKeyFile())
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
	
	return tmtNode.NewNode(
		cfg,
		privval.LoadOrGenFilePV(cfg.PrivValidatorFile()),
		nodeKey,
		proxy.DefaultClientCreator(cfg.ProxyApp, cfg.ABCI, cfg.DBDir()),
		tmtNode.DefaultGenesisDocProviderFunc(cfg),
		tmtNode.DefaultDBProvider,
		tmtNode.DefaultMetricsProvider(cfg.Instrumentation),
		logger,
	)
}

