package node

import (
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/log"
)

type Node struct {
	dbNode *database.Node
	consensusNode *consensus.Node
	logger log.Logger
}

// makeFullNode creates a full go-database node
func NewNode(cfg *Config) (*Node, error) {
	
	dbNode, err := database.NewNode(&cfg.dbCfg)
	if err != nil {
		return nil, err
	}
	
	consensusNode, err := consensus.NewNode(&cfg.consensusCfg)
	if err != nil {
		return nil, err
	}
	
	nodeLogger := log.NewLogger()
	nodeLogger.With("module", "node")
	return &Node {
		dbNode,
		consensusNode,
		nodeLogger,
	}, nil
}


// Start starts base node and stop p2p server
func (n *Node) Start() error {
	// Start database node
	uriClient := n.consensusNode.NewURIClient()
	err := n.dbNode.Start(uriClient)
	if err != nil {
		return err
	}

	n.consensusNode.Start()
	if err != nil {
		return err
	}
	return nil
}


func (n *Node) Stop() error {
	err := n.dbNode.Stop()
	if err != nil {
		return err
	}
	
	n.consensusNode.Stop()
	if err != nil {
		return err
	}

	return nil
}

