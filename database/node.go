package database

import (
	"os"
	"fmt"

	ethNode "github.com/ethereum/go-ethereum/node"
	ethRpc "github.com/ethereum/go-ethereum/rpc"

	rpcClient "github.com/tendermint/tendermint/rpc/lib/client"

	"github.com/lightstreams-network/lightchain/log"
	"google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Node is the main object.
type Node struct {
	ethereum  *ethNode.Node
	database  *Database
	rpcClient *ethRpc.Client
	cfg       *Config
	logger    log.Logger
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

	n.logger.Debug("Register wallet event handlers to open and auto-derive wallets")
	go registerEventHandlers(n)

	// Register NewNode ABCI Application Service
	n.logger.Debug("Binding ethereum events to rpc client...")
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
	n.logger.Debug("Starting ethereum service...")
	if err := n.ethereum.Service(&ethDb); err != nil {
		n.logger.Error(fmt.Errorf("database service not running: %v", err).Error())
		return err
	}

	n.rpcClient, err = n.ethereum.Attach()
	if err != nil {
		n.logger.Error(fmt.Errorf("failed to attach to the inproc geth: %v", err).Error())
		os.Exit(1)
	}

	n.logger.Info("Started database engine...")
	return nil
}

func (n *Node) Stop() error {
	n.logger.Debug("Stopping ethereum service")
	if err := n.ethereum.Stop(); err != nil {
		return err
	}

	n.logger.Debug("Stopping database service")
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

func registerEventHandlers(n *Node) error {
	events := make(chan accounts.WalletEvent, 16)
	n.ethereum.AccountManager().Subscribe(events)

	// Create an chain state reader for self-derivation
	client, err := n.ethereum.Attach() // Ethereum RPC client
	if err != nil {
		n.logger.Error(fmt.Errorf("Failed to attach to self: %v", err).Error())
		return err
	}
	stateReader := ethclient.NewClient(client)

	// Open and self derive any wallets already attached
	for _, wallet := range n.ethereum.AccountManager().Wallets() {
		if err := wallet.Open(""); err != nil {
			n.logger.Warn("Failed to open wallet", "url", wallet.URL(), "err", err)
		} else {
			wallet.SelfDerive(accounts.DefaultBaseDerivationPath, stateReader)
		}
	}

	// Listen for wallet event till termination
	for event := range events {
		switch event.Kind {
		case accounts.WalletArrived:
			if err := event.Wallet.Open(""); err != nil {
				n.logger.Warn("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
			}
		case accounts.WalletOpened:
			status, _ := event.Wallet.Status()
			n.logger.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

			if event.Wallet.URL().Scheme == "ledger" {
				event.Wallet.SelfDerive(accounts.DefaultLedgerBaseDerivationPath, stateReader)
			} else {
				event.Wallet.SelfDerive(accounts.DefaultBaseDerivationPath, stateReader)
			}

		case accounts.WalletDropped:
			n.logger.Info("Old wallet dropped", "url", event.Wallet.URL())
			event.Wallet.Close()
		}
	}

	return nil
}
