package consensus

import (
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/tendermint/tendermint/proxy"
	"github.com/lightstreams-network/lightchain/database"

	tmtNode "github.com/tendermint/tendermint/node"
	tmtServer "github.com/tendermint/tendermint/abci/server"
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

func (n *Node) Start(rpcClient *ethRpc.Client, db *database.Database) error {
	if err := n.tendermint.Start(); err != nil {
		return err
	}

	tendermintABCI, err := NewTendermintABCI(db, rpcClient, n.logger)
	if err != nil {
		return err
	}

	err = tendermintABCI.InitEthState()
	if err != nil {
		return err
	}

	abciSrv, err := tmtServer.NewServer(n.cfg.tendermintCfg.ProxyApp, n.cfg.proxyProtocol, tendermintABCI)
	if err != nil {
		return err
	}

	abciSrv.SetLogger(log.NewLogger().With("module", "node-server"))

	if err := abciSrv.Start(); err != nil {
		return err
	}

	n.logger.Info("Consensus node started", "nodeInfo", n.tendermint.Switch().NodeInfo())

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