package database

import (
	"github.com/ethereum/go-ethereum/node"
	
)

// Node is the main object.
type Node struct {
	node.Node
}

// NewNode creates a new node.
func NewNode(conf *node.Config) (*Node, error) {
	stack, err := node.New(conf)
	if err != nil {
		return nil, err
	}

	return &Node{*stack}, nil // nolint: vet
}

// Start starts base node and stop p2p server
func (n *Node) Start() error {
	// start p2p server
	err := n.Node.Start()
	if err != nil {
		return err
	}

	// Stop it Eth.p2p server
	n.Node.Server().Stop()

	return nil
}


func (n *Node) Stop() error {
	err := n.Node.Stop()
	if err != nil {
		return err
	}

	return nil
}
