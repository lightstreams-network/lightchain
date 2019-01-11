package database

import (
	"google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"fmt"
	"os"
	
	ethNode "github.com/ethereum/go-ethereum/node"
	
	rpcClient "github.com/tendermint/tendermint/rpc/lib/client"
	
	"github.com/lightstreams-network/lightchain/log"
)

// Node is the main object.
type Node struct {
	ethereum *ethNode.Node
	cfg      *Config
	logger	log.Logger
}

// NewNode creates a new node.
func NewNode(cfg *Config) (*Node, error) {
	dbLogger := log.NewLogger()
	dbLogger.With("module", "database")
	stack, err := ethNode.New(&cfg.GethConfig.Node)
	if err != nil {
		return nil, err
	}

	return &Node{
		stack,
		cfg,
		dbLogger,
	}, nil // nolint: vet
}

// Start starts base node and stop p2p server
func (n *Node) Start(uriClient *rpcClient.URIClient) error {
	// start p2p server
	n.logger.Info("Starting ethereum node")
	err := n.ethereum.Start()
	if err != nil {
		return err
	}

	// Stop it Eth.p2p server
	n.ethereum.Server().Stop()
	
	// Register NewNode ABCI Application Service
	if err := n.ethereum.Register(func(ctx *ethNode.ServiceContext) (ethNode.Service, error) {
		n.logger.Info(fmt.Sprintf("registering ABCI application service"))
		return New(ctx, &n.cfg.GethConfig.Eth, uriClient)
	}); err != nil {
		return err
	}
	
	// Fetch the registered service of this type
	var ethDb *database.Database
	if err := n.ethereum.Service(&ethDb); err != nil {
		n.logger.Error(fmt.Errorf("database ethBackend service not running: %v", err).Error())
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
