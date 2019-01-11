package database

import (
	ethNode "github.com/ethereum/go-ethereum/node"
	rpcClient "github.com/tendermint/tendermint/rpc/lib/client"
	
)

// Node is the main object.
type Node struct {
	ethereum *ethNode.Node
	cfg      *Config
}

// NewNode creates a new node.
func NewNode(cfg *Config) (*Node, error) {
	stack, err := ethNode.New(&cfg.GethConfig.Node)
	if err != nil {
		return nil, err
	}

	return &Node{
		stack,
		cfg,
	}, nil // nolint: vet
}

// Start starts base node and stop p2p server
func (n *Node) Start(uriClient *rpcClient.URIClient) error {
	// start p2p server
	err := n.ethereum.Start()
	if err != nil {
		return err
	}

	// Stop it Eth.p2p server
	n.ethereum.Server().Stop()
	
	// Register NewNode ABCI Application Service
	if err := n.ethereum.Register(func(ctx *ethNode.ServiceContext) (ethNode.Service, error) {
		return New(ctx, &n.cfg.GethConfig.Eth, uriClient)
	}); err != nil {
		return err
	}

	return nil
}


func (n *Node) Stop() error {
	err := n.ethereum.Stop()
	if err != nil {
		return err
	}

	return nil
}
