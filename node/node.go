package node

import (
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
)

type Node struct {
	dbNode *database.Node
	consensusNode *consensus.Node
}


// Start starts base node and stop p2p server
func (n *Node) Start() error {
	// start p2p server
	err := n.dbNode.Start()
	if err != nil {
		return err
	}

	// Stop it Eth.p2p server
	n.dbNode.Server().Stop()

	return nil
}


func (n *Node) Stop() error {
	err := n.dbNode.Stop()
	if err != nil {
		return err
	}

	return nil
}
