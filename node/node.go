package node

import (
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/prometheus"
	"github.com/lightstreams-network/lightchain/log"
	conAPI "github.com/lightstreams-network/lightchain/consensus/api"
	tmtLog "github.com/tendermint/tendermint/libs/log"
)

type Node struct {
	cfg            *Config
	dbNode         *database.Node
	consensusNode  *consensus.Node
	prometheusNode *prometheus.Node
	logger         tmtLog.Logger
}

func NewNode(cfg *Config) (*Node, error) {
	logger := log.NewLogger().With("engine", "node")

	logger.Debug("Creating new prometheus node instance from config...")
	prometheusNode := prometheus.NewNode(cfg.prometheusCfg)

	logger.Debug("Creating new consensus node instance from config...")
	consensusNode, err := consensus.NewNode(&cfg.consensusCfg, prometheusNode.Registry())
	if err != nil {
		return nil, err
	}

	// Todo: #90 Should be refactored, creating a new instance of a struct SHOULDN'T do any FS changes
	conRPCAPI := conAPI.NewRPCApi(cfg.consensusCfg.RPCListenPort())
	logger.Debug("Creating new database node instance from config and creates a keystore dir...")
	dbNode, err := database.NewNode(&cfg.dbCfg, conRPCAPI, prometheusNode.Registry())
	if err != nil {
		return nil, err
	}

	return &Node{
		cfg,
		dbNode,
		consensusNode,
		prometheusNode,
		logger,
	}, nil
}

func (n *Node) Start() error {
	n.logger.Info("Starting database engine...")
	if err := n.dbNode.Start(); err != nil {
		return err
	}

	n.logger.Info("Starting consensus engine...")
	if err := n.consensusNode.Start(n.dbNode.RpcClient(), n.dbNode.Database()); err != nil {
		return err
	}

	if err := n.prometheusNode.Start(); err != nil {
		return err
	}

	return nil
}

func (n *Node) Stop() error {
	// IMPORTANT: We need to close consensus first so that node stops receiving new blocks
	// before database is closed
	n.logger.Info("Stopping consensus engine...")
	if err := n.consensusNode.Stop(); err != nil {
		return err
	}
	n.logger.Info("Consensus node stopped")

	n.logger.Info("Stopping database engine...")
	if err := n.dbNode.Stop(); err != nil {
		return err
	}
	n.logger.Info("Database node stopped")

	if err := n.prometheusNode.Stop(); err != nil {
		return err
	}

	return nil
}