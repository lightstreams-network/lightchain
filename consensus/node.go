package consensus

import (
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/tendermint/tendermint/proxy"

	tmtNode "github.com/tendermint/tendermint/node"
	rpcTypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcClient "github.com/tendermint/tendermint/rpc/lib/client"
	ethRpc "github.com/ethereum/go-ethereum/rpc"
)

type Node struct {
	tendermint *tmtNode.Node
	cfg        *Config
	logger     log.Logger
} 

func NewNode(cfg *Config) (*Node, error) {
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.tendermintCfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	consensusLogger := log.NewLogger()
	consensusLogger.With("module", "consensus")

	tendermint, err := tmtNode.NewNode(
		cfg.tendermintCfg,
		privval.LoadOrGenFilePV(cfg.tendermintCfg.PrivValidatorFile()),
		nodeKey,
		proxy.DefaultClientCreator(cfg.tendermintCfg.ProxyApp, cfg.tendermintCfg.ABCI, cfg.tendermintCfg.DBDir()),
		tmtNode.DefaultGenesisDocProviderFunc(cfg.tendermintCfg),
		tmtNode.DefaultDBProvider,
		tmtNode.DefaultMetricsProvider(cfg.tendermintCfg.Instrumentation),
		consensusLogger,
	)

	return &Node{
		tendermint,
		cfg,
		consensusLogger,
	}, nil
}

func (n *Node) Start(rpcClient *ethRpc.Client) error {
	if err := n.tendermint.Start(); err != nil {
		return err
	}

	n.logger.Info("Started node", "nodeInfo", n.tendermint.Switch().NodeInfo())

	return nil
}

func (n *Node) Stop() error {
	return nil
}

func (n *Node) NewURIClient() *rpcClient.URIClient {
	client := rpcClient.NewURIClient(n.cfg.tendermintCfg.RPC.ListenAddress)
	rpcTypes.RegisterAmino(client.Codec())

	return client
}