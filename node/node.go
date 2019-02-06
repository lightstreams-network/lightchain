package node

import (
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/log"
	conAPI "github.com/lightstreams-network/lightchain/consensus/api"
	tmtLog "github.com/tendermint/tendermint/libs/log"
)

type Node struct {
	dbNode        *database.Node
	consensusNode *consensus.Node
	logger        tmtLog.Logger
}

// makeFullNode creates a full go-database node
func NewNode(cfg *Config) (*Node, error) {
	logger := log.NewLogger().With("engine", "node")
	logger.Debug("Initializing consensus node...")
	consensusNode, err := consensus.NewNode(&cfg.consensusCfg)
	if err != nil {
		return nil, err
	}

	conRPCAPI := conAPI.NewRPCApi(cfg.consensusCfg.RPCListenPort())

	logger.Debug("Initializing database node...")
	dbNode, err := database.NewNode(&cfg.dbCfg, conRPCAPI)
	if err != nil {
		return nil, err
	}

	return &Node{
		dbNode,
		consensusNode,
		logger,
	}, nil
}

// Start starts base node and stop p2p server
func (n *Node) Start() error {
	// Start database node
	n.logger.Info("Starting database engine...")
	if err := n.dbNode.Start(); err != nil {
		return err
	}
	
	n.logger.Info("Starting consensus engine...")
	if err := n.consensusNode.Start(n.dbNode.RpcClient(), n.dbNode.Database()); err != nil {
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
	return nil
}
