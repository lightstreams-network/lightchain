package node

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"path/filepath"
	"io/ioutil"
	
	"github.com/ethereum/go-ethereum/accounts"
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	ethLog "github.com/ethereum/go-ethereum/log"
	
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	
	rpcTypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcClient "github.com/tendermint/tendermint/rpc/lib/client"
	
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/utils"
)

func InitNode(ctx *cli.Context) error {
	homeDir := utils.MakeHomeDir(ctx)
	ethLog.Info("Initializing HomeDir", "dir", homeDir)
	dataDir := utils.MakeDataDir(ctx)

	// Lightstreams configs
	lsCfgPath := filepath.Join(dataDir, consensus.ConfigFilename)
	err := ioutil.WriteFile(lsCfgPath, consensus.ReadDefaultConfig(), 0666)
	if err != nil {
		ethUtils.Fatalf("Config err: %v", err)
	} else {
		ethLog.Info(fmt.Sprintf("successfully copied config into: %s", lsCfgPath))
	}

	return nil
}

// startNode copies the logic from go-database (go-database/cmd/geth/main.go)
func StartNode(ctx *cli.Context, stack *database.Node) {
	if err := stack.Start(); err != nil {
		ethUtils.Fatalf("Error starting protocol stack: %v", err)
	}
	
	// Register wallet event handlers to open and auto-derive wallets
	events := make(chan accounts.WalletEvent, 16)
	stack.AccountManager().Subscribe(events)

	go func() {
		// Create an chain state reader for self-derivation
		client, err := stack.Attach() // Ethereum RPC client
		if err != nil {
			ethUtils.Fatalf("Failed to attach to self: %v", err)
		}
		stateReader := ethclient.NewClient(client)

		// Open and self derive any wallets already attached
		for _, wallet := range stack.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				log.Warn("Failed to open wallet", "url", wallet.URL(), "err", err)
			} else {
				wallet.SelfDerive(accounts.DefaultBaseDerivationPath, stateReader)
			}
		}

		// Listen for wallet event till termination
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					log.Warn("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				log.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

				if event.Wallet.URL().Scheme == "ledger" {
					event.Wallet.SelfDerive(accounts.DefaultLedgerBaseDerivationPath, stateReader)
				} else {
					event.Wallet.SelfDerive(accounts.DefaultBaseDerivationPath, stateReader)
				}

			case accounts.WalletDropped:
				log.Info("Old wallet dropped", "url", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()
}

// makeFullNode creates a full go-database node
func CreateNode(ctx *cli.Context) *database.Node {
	tmtCfg := consensus.MakeTendermintConfig(ctx)
	
	// Step 1: Setup the go-database node and start it
	tendermintLAddr := fmt.Sprintf("tcp://%s:%d", "127.0.0.1", tmtCfg.RpcListenPort)
	stack, cfg := makeConfigNode(ctx)

	// Register New ABCI Application Service
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		client := rpcClient.NewURIClient(tendermintLAddr) // tendermint RPC client
		rpcTypes.RegisterAmino(client.Codec())
		return database.NewBackend(ctx, &cfg.Eth, client)
	}); err != nil {
		ethUtils.Fatalf("Failed to register the ABCI application service: %v", err)
	}

	return stack
}

func makeConfigNode(ctx *cli.Context) (*database.Node, database.GethConfig) {
	cfg := database.GethConfig{
		Eth:  eth.DefaultConfig,
		Node: database.DefaultNodeConfig(),
	}

	ethUtils.SetNodeConfig(ctx, &cfg.Node)
	database.SetNodeConfig(ctx, &cfg.Node)
 	stack, err := database.New(&cfg.Node)
	if err != nil {
		ethUtils.Fatalf("Failed to create the protocol stack: %v", err)
	}

	ethUtils.SetEthConfig(ctx, &stack.Node, &cfg.Eth)
	database.SetEthConfig(&cfg.Eth)

	return stack, cfg
}
