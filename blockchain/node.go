package blockchain

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"strings"
	"path/filepath"
	"os"
	"io/ioutil"
	
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	ethLog "github.com/ethereum/go-ethereum/log"
	
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/core"
	
	rpcTypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcClient "github.com/tendermint/tendermint/rpc/lib/client"
	
	"github.com/lightstreams-network/lightchain/ethereum"
	"github.com/lightstreams-network/lightchain/config"
	"github.com/lightstreams-network/lightchain/utils"
)

func InitNode(ctx *cli.Context) error {
	// Step 1: Init chain within --datadir by read genesis
	homeDir := config.MakeDataDir(ctx)
	ethLog.Info("Initializing HomeDir", "dir", homeDir)
	dataDir := filepath.Join(homeDir, config.ConfigFolderName)
	
	chainDb, err := ethdb.NewLDBDatabase(filepath.Join(dataDir, "chaindata"), 0, 0)
	if err != nil {
		ethUtils.Fatalf("could not open database: %v", err)
	}

	keystoreDir := filepath.Join(homeDir, "keystore")
	if err := os.MkdirAll(keystoreDir, os.ModePerm); err != nil {
		ethUtils.Fatalf("mkdirAll keyStoreDir: %v", err)
	}

	keystoreCfg, err := config.ReadDefaultKeystore()
	if err != nil {
		ethUtils.Fatalf("could not open read keystore: %v", err)
	}

	for filename, content := range keystoreCfg {
		storeFileName := filepath.Join(keystoreDir, filename)
		f, err := os.Create(storeFileName)
		if err != nil {
			ethLog.Error("Cannot create file", storeFileName, err)
			continue
		}
		if _, err := f.Write(content); err != nil {
			ethLog.Error("write content %q err: %v", storeFileName, err)
		}
		if err := f.Close(); err != nil {
			return err
		}

		ethLog.Info("Successfully wrote keystore files", "keystore", storeFileName)
	}

	genesisPath := config.MakeGenesisPath(ctx)
	ethLog.Info("Trying to reading genesis", "dir", genesisPath)
	genesisBlob, err := config.ReadGenesisPath(genesisPath)
	if err != nil {
		ethLog.Warn("Error reading genesisPath", err)
		genesisBlob, err = config.ReadDefaultGenesis()
		if err != nil {
			ethUtils.Fatalf("genesis read error: %v", err)
		}
	}
	genesis, err := utils.ParseBlobGenesis(genesisBlob)
	if err != nil {
		ethUtils.Fatalf("genesisJSON err: %v", err)
	}

	genesisFileName := filepath.Join(homeDir, "genesis.json")
	f, err := os.Create(genesisFileName)
	if _, err := f.Write(genesisBlob); err != nil {
		ethLog.Error("write content %q err: %v", genesisFileName, err)
	}

	ethLog.Info("Using genesis block", "block", genesis)

	_, hash, err := core.SetupGenesisBlock(chainDb, genesis)
	if err != nil {
		ethUtils.Fatalf("failed to write genesis block: %v", err)
	}

	ethLog.Info("Successfully wrote genesis block and/or chain rule set", "hash", hash)

	// Lightstreams configs
	lsCfgPath := filepath.Join(homeDir, config.ConfigFolderName, config.ConfigFilename)
	err = ioutil.WriteFile(lsCfgPath, config.ReadDefaultLsConfigBlob(), 0666)
	if err != nil {
		ethUtils.Fatalf("Config err: %v", err)
	} else {
		ethLog.Info(fmt.Sprintf("successfully copied LS config into: %s", lsCfgPath))
	}

	return nil
}

// startNode copies the logic from go-ethereum (go-ethereum/cmd/geth/main.go)
func StartNode(ctx *cli.Context, stack *ethereum.Node) {
	if err := stack.Start(); err != nil {
		ethUtils.Fatalf("Error starting protocol stack: %v", err)
	}

	// Unlock any account specifically requested
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

	passwords := ethUtils.MakePasswordList(ctx)
	unlocks := strings.Split(ctx.GlobalString(ethUtils.UnlockedAccountFlag.Name), ",")
	for i, account := range unlocks {
		if trimmed := strings.TrimSpace(account); trimmed != "" {
			unlockAccount(ks, trimmed, i, passwords)
		}
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

// makeFullNode creates a full go-ethereum node
func CreateNode(tendermintLAddr string, ctx *cli.Context) *ethereum.Node {
	stack, cfg := makeConfigNode(ctx)

	// Register New ABCI Application Service
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		client := rpcClient.NewURIClient(tendermintLAddr) // tendermint RPC client
		rpcTypes.RegisterAmino(client.Codec())
		return ethereum.NewBackend(ctx, &cfg.Eth, client)
	}); err != nil {
		ethUtils.Fatalf("Failed to register the ABCI application service: %v", err)
	}

	return stack
}

func makeConfigNode(ctx *cli.Context) (*ethereum.Node, config.GethConfig) {
	cfg := config.GethConfig{
		Eth:  eth.DefaultConfig,
		Node: config.DefaultNodeConfig(),
	}

	config.SetLightchainNodeDefaultConfig(&cfg.Node)
	ethUtils.SetNodeConfig(ctx, &cfg.Node)
	stack, err := ethereum.New(&cfg.Node)
	if err != nil {
		ethUtils.Fatalf("Failed to create the protocol stack: %v", err)
	}

	config.SetLightchainEthDefaultConfig(&cfg.Eth)
	ethUtils.SetEthConfig(ctx, &stack.Node, &cfg.Eth)

	return stack, cfg
}

// tries unlocking the specified account a few times.
func unlockAccount(ks *keystore.KeyStore,
	address string,
	i int,
	passwords []string,
) (accounts.Account, string) {

	account, err := ethUtils.MakeAddress(ks, address)
	if err != nil {
		ethUtils.Fatalf("Could not list accounts: %v", err)
	}
	for trials := 0; trials < 3; trials++ {
		prompt := fmt.Sprintf("Unlocking account %s | Attempt %d/%d", address, trials+1, 3)
		password := getPassPhrase(prompt, false, i, passwords)
		err = ks.Unlock(account, password)
		if err == nil {
			log.Info("Unlocked account", "address", account.Address.Hex())
			return account, password
		}
		if err, ok := err.(*keystore.AmbiguousAddrError); ok {
			log.Info("Unlocked account", "address", account.Address.Hex())
			return ambiguousAddrRecovery(ks, err, password), password
		}
		if err != keystore.ErrDecrypt {
			// No need to prompt again if the error is not decryption-related.
			break
		}
	}
	// All trials expended to unlock account, bail out
	ethUtils.Fatalf("Failed to unlock account %s (%v)", address, err)

	return accounts.Account{}, ""
}

// getPassPhrase retrieves the passwor associated with an account, either fetched
// from a list of preloaded passphrases, or requested interactively from the user.
// nolint: unparam
func getPassPhrase(prompt string,
	confirmation bool,
	i int,
	passwords []string,
) string {
	// If a list of passwords was supplied, retrieve from them
	if len(passwords) > 0 {
		if i < len(passwords) {
			return passwords[i]
		}
		return passwords[len(passwords)-1]
	}
	// Otherwise prompt the user for the password
	if prompt != "" {
		fmt.Println(prompt)
	}
	password, err := console.Stdin.PromptPassword("Passphrase: ")
	if err != nil {
		ethUtils.Fatalf("Failed to read passphrase: %v", err)
	}
	if confirmation {
		confirm, err := console.Stdin.PromptPassword("Repeat passphrase: ")
		if err != nil {
			ethUtils.Fatalf("Failed to read passphrase confirmation: %v", err)
		}
		if password != confirm {
			ethUtils.Fatalf("Passphrases do not match")
		}
	}
	return password
}

func ambiguousAddrRecovery(ks *keystore.KeyStore, err *keystore.AmbiguousAddrError,
	auth string) accounts.Account {

	fmt.Printf("Multiple key files exist for address %x:\n", err.Addr)
	for _, a := range err.Matches {
		fmt.Println("  ", a.URL)
	}
	fmt.Println("Testing your passphrase against all of them...")
	var match *accounts.Account
	for _, a := range err.Matches {
		if err := ks.Unlock(a, auth); err == nil {
			match = &a
			break
		}
	}
	if match == nil {
		ethUtils.Fatalf("None of the listed files could be unlocked.")
	}
	fmt.Printf("Your passphrase unlocked %s\n", match.URL)
	fmt.Println("In order to avoid this warning, remove the following duplicate key files:")
	for _, a := range err.Matches {
		if a != *match {
			fmt.Println("  ", a.URL)
		}
	}
	return *match
}
