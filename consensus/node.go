package consensus

import (
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/tendermint/tendermint/proxy"
	"github.com/lightstreams-network/lightchain/database"

	tmtNode "github.com/tendermint/tendermint/node"
	tmtP2P "github.com/tendermint/tendermint/p2p"
	tmtCommon "github.com/tendermint/tendermint/libs/common"
	tmtServer "github.com/tendermint/tendermint/abci/server"
	rpcTypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcClient "github.com/tendermint/tendermint/rpc/lib/client"
	ethRpc "github.com/ethereum/go-ethereum/rpc"
	"fmt"
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
	
	nodeKeyFile := cfg.tendermintCfg.NodeKeyFile()
	if ! tmtCommon.FileExists(nodeKeyFile) {
		return nil, fmt.Errorf("Tendermint key file does not exists")
	}
	
	nodeKey, err := p2p.LoadOrGenNodeKey(nodeKeyFile)
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

func (n *Node) Start(ethRPCClient *ethRpc.Client, db *database.Database) error {
	n.logger.Debug("Creating tendermint ABCI application...")
	tendermintABCI, err := NewTendermintABCI(db, ethRPCClient, n.logger)
	if err != nil {
		return err
	}

	n.logger.Debug("Initializing consensus state...")
	err = tendermintABCI.InitEthState()
	if err != nil {
		return err
	}

	n.logger.Debug("Creating tendermint server...")
	proxyLAddr := fmt.Sprintf("tcp://0.0.0.0:%d", n.cfg.proxyListenPort)
	abciSrv, err := tmtServer.NewServer(proxyLAddr, n.cfg.proxyProtocol, tendermintABCI)
	if err != nil {
		return err
	}

	n.logger.Debug("Starting ABCI application server...")
	abciSrv.SetLogger(log.NewLogger().With("module", "node-server"))
	if err := abciSrv.Start(); err != nil {
		return err
	}
	n.logger.Debug("Started ABCI application server")

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
	
	
	n.logger.Info("Starting tendermint node...")
	if err := n.tendermint.Start(); err != nil {
		return err
	}
	n.logger.Info("Started tendermint node")

	n.logger.Info("Consensus node started", "nodeInfo", n.tendermint.Switch().NodeInfo())
	return nil
}

func (n *Node) Stop() error {
	return nil
}

func (n *Node) NewURIClient() *rpcClient.URIClient {
	tendermintLAddr := fmt.Sprintf("tcp://127.0.0.1:%d", n.cfg.rpcListenPort)
	client := rpcClient.NewURIClient(tendermintLAddr)
	rpcTypes.RegisterAmino(client.Codec())
	return client
}
