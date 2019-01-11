package consensus

import (
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/tendermint/tendermint/proxy"
	"github.com/lightstreams-network/lightchain/database"

	tmtNode "github.com/tendermint/tendermint/node"
	tmtP2P "github.com/tendermint/tendermint/p2p"
	tmtServer "github.com/tendermint/tendermint/abci/server"
	rpcTypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcClient "github.com/tendermint/tendermint/rpc/lib/client"
	ethRpc "github.com/ethereum/go-ethereum/rpc"
)

type Node struct {
	tendermint *tmtNode.Node
	nodeKey	   *tmtP2P.NodeKey
	cfg        *Config
	logger     log.Logger
}

func NewNode(cfg *Config) (*Node, error) {
	consensusLogger := log.NewLogger()
	consensusLogger.With("module", "consensus")
	
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.tendermintCfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	return &Node{
		nil,
		nodeKey,
		cfg,
		consensusLogger,
	}, nil
}

func (n *Node) Start(rpcClient *ethRpc.Client, db *database.Database) error {
	
	n.logger.Info("Creating tendermint node...")
	tendermint, err := tmtNode.NewNode(
		n.cfg.tendermintCfg,
		privval.LoadOrGenFilePV(n.cfg.tendermintCfg.PrivValidatorFile()),
		n.nodeKey,
		proxy.DefaultClientCreator(n.cfg.tendermintCfg.ProxyApp, n.cfg.tendermintCfg.ABCI, n.cfg.tendermintCfg.DBDir()),
		tmtNode.DefaultGenesisDocProviderFunc(n.cfg.tendermintCfg),
		tmtNode.DefaultDBProvider,
		tmtNode.DefaultMetricsProvider(n.cfg.tendermintCfg.Instrumentation),
		n.logger,
	)

	if err != nil {
		return err
	} else {
		n.tendermint = tendermint
	}
	
	
	n.logger.Info("Starting consensus engine...")
	if err := n.tendermint.Start(); err != nil {
		return err
	}

	n.logger.Debug("Creating tendermint ABCI application...")
	tendermintABCI, err := NewTendermintABCI(db, rpcClient, n.logger)
	if err != nil {
		return err
	}

	n.logger.Debug("Initializing consensus state...")
	err = tendermintABCI.InitEthState()
	if err != nil {
		return err
	}

	n.logger.Debug("Creating tendermint server...")
	abciSrv, err := tmtServer.NewServer(n.cfg.tendermintCfg.ProxyApp, n.cfg.proxyProtocol, tendermintABCI)
	if err != nil {
		return err
	}

	abciSrv.SetLogger(log.NewLogger().With("module", "node-server"))
	if err := abciSrv.Start(); err != nil {
		return err
	}

	n.logger.Info("Consensus node started", "nodeInfo", n.tendermint.Switch().NodeInfo())

	n.logger.Info("Started consensus node...")
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
