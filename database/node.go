package database

import (
	"os"
	"fmt"
	
	ethNode "github.com/ethereum/go-ethereum/node"
	ethRpc "github.com/ethereum/go-ethereum/rpc"
	
	rpcClient "github.com/tendermint/tendermint/rpc/lib/client"
	
	"github.com/lightstreams-network/lightchain/log"
	"google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

// Node is the main object.
type Node struct {
	ethereum *ethNode.Node
	database *Database
	rpcClient *ethRpc.Client
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
		nil,
		nil,
		cfg,
		dbLogger,
	}, nil // nolint: vet
}

// Start starts base node and stop p2p server
func (n *Node) Start(uriClient *rpcClient.URIClient) error {
	// start p2p server
	n.logger.Debug("Starting ethereum node")
	err := n.ethereum.Start()
	if err != nil {
		return err
	}

	// Stop it Eth.p2p server
	n.logger.Debug("Stopping p2p communication of ethereum node")
	n.ethereum.Server().Stop()
	
	// Register NewNode ABCI Application Service
	if err := n.ethereum.Register(func(ctx *ethNode.ServiceContext) (ethNode.Service, error) {
		n.logger.Info(fmt.Sprintf("registering ABCI application service"))
		n.database, err = NewDatabase(ctx, &n.cfg.GethConfig.Eth, uriClient)
		if err != nil {
			return nil, err
		}
		return n.database, nil
	}); err != nil {
		return err
	}
	
	// Fetch the registered service of this type
	var ethDb *database.Database
	if err := n.ethereum.Service(&ethDb); err != nil {
		n.logger.Error(fmt.Errorf("database service not running: %v", err).Error())
		return err
	}
	
	n.rpcClient, err = n.ethereum.Attach()
	if err != nil {
		n.logger.Error(fmt.Errorf("failed to attach to the inproc geth: %v", err).Error())
		os.Exit(1)
	}

	return nil
}


func (n *Node) Stop() error {
	if err := n.ethereum.Stop(); err != nil {
		return err
	}
	
	if err := n.database.Stop(); err != nil {
		return err
	}

	return nil
}

func (n *Node) RpcClient() *ethRpc.Client {
	return n.rpcClient
}

func (n *Node) Database() *Database {
	return n.database
}
