package main

import (
	"fmt"
	"os"

	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	mintApp "github.com/lightstreams-network/lightchain/app"
	emtUtils "github.com/lightstreams-network/lightchain/cmd/utils"
	"github.com/lightstreams-network/lightchain/ethereum"
	"github.com/lightstreams-network/lightchain/version"
	"github.com/tendermint/tendermint/abci/server"
	mintLog "github.com/tendermint/tendermint/libs/log"
	"github.com/lightstreams-network/lightchain/core/transaction"
	lsConfig "github.com/lightstreams-network/lightchain/core/config"
	cmn "github.com/tendermint/tmlibs/common"
	"gopkg.in/urfave/cli.v1"
	"path/filepath"
	"io/ioutil"
)

var (
	// The app that holds all commands and flags.
	app = ethUtils.NewApp(version.Version, "the lightchain command line interface")
)

func BeforeCmd(ctx *cli.Context) error {
	if err := emtUtils.Setup(ctx); err != nil {
		return err
	}
	//ethUtils.SetupNetwork(ctx)
	return nil
}

func LightchainNodeCmd(ctx *cli.Context) {
	fmt.Println("Start Node")
	// Step 1: Setup the go-ethereum node and start it
	node := emtUtils.MakeFullNode(ctx)
	emtUtils.StartNode(ctx, node)

	// Setup the ABCI server and start it
	addr := ctx.GlobalString(emtUtils.ABCIAddrFlag.Name)
	abci := ctx.GlobalString(emtUtils.ABCIProtocolFlag.Name)

	// Fetch the registered service of this type
	var ethBackend *ethereum.Backend
	if err := node.Service(&ethBackend); err != nil {
		ethUtils.Fatalf("ethereum ethBackend service not running: %v", err)
	}

	// In-proc RPC connection so ABCI.Query can be forwarded over the ethereum rpc
	rpcClient, err := node.Attach()
	if err != nil {
		ethUtils.Fatalf("Failed to attach to the inproc geth: %v", err)
	}
	
	txManager, err := transaction.NewManager(buildLsConfigPath(ctx))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create the ABCI app
	// Create the application - in memory or persisted to disk
	ethApp, err := mintApp.CreateLightchainApplication(ethBackend, rpcClient, nil, txManager)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ethApp.SetLogger(emtUtils.LightchainLogger().With("module", "lightchain"))

	// Start the app on the ABCI server listener
	appSrv, err := server.NewServer(addr, abci, ethApp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger := mintLog.NewTMLogger(mintLog.NewSyncWriter(os.Stdout))
	appSrv.SetLogger(logger.With("module", "abci-server"))

	if err := appSrv.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmn.TrapSignal(func() {
		appSrv.Stop()
	})
}

func kvStoreCmd(ctx *cli.Context) {
	logger := mintLog.NewTMLogger(mintLog.NewSyncWriter(os.Stdout))
	// Setup the ABCI server and start it
	flagAddress := ctx.GlobalString(emtUtils.ABCIAddrFlag.Name)
	flagAbci := ctx.GlobalString(emtUtils.ABCIProtocolFlag.Name)

	// Create the application - in memory
	kvApp := mintApp.NewKVStoreApplication()
	kvApp.SetLogger(emtUtils.LightchainLogger().With("module", "kwstore"))

	// Start the app on the ABCI server listener
	srv, err := server.NewServer(flagAddress, flagAbci, kvApp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	srv.SetLogger(logger.With("module", "abci-server"))

	if err := srv.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})
}

func VersionCmd(ctx *cli.Context) error {
	fmt.Println("lightchain: ", version.Version)
	fmt.Println("go-ethereum: ", params.Version)
	return nil
}

func InitCmd(ctx *cli.Context) error {
	// Step 1: Init chain within --datadir by read genesis
	lightchainDataDir := emtUtils.MakeDataDir(ctx)
	ethLog.Info("Initializing HomeDir", "dir", lightchainDataDir)

	chainDb, err := ethdb.NewLDBDatabase(filepath.Join(lightchainDataDir,
		"lightchain/chaindata"), 0, 0)
	if err != nil {
		ethUtils.Fatalf("could not open database: %v", err)
	}

	keystoreDir := filepath.Join(lightchainDataDir, "keystore")
	if err := os.MkdirAll(keystoreDir, 0666); err != nil {
		ethUtils.Fatalf("mkdirAll keyStoreDir: %v", err)
	}

	keystore, err := emtUtils.ReadDefaultKeystore()
	if err != nil {
		ethUtils.Fatalf("could not open read keystore: %v", err)
	}

	for filename, content := range keystore {
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

	genesisPath := emtUtils.MakeGenesisPath(ctx)
	ethLog.Info("Trying to reading genesis", "dir", genesisPath)
	genesisBlob, err := emtUtils.ReadGenesisPath(genesisPath)
	if err != nil {
		ethLog.Warn("Error reading genesisPath", err)
		genesisBlob, err = emtUtils.ReadDefaultGenesis()
		if err != nil {
			ethUtils.Fatalf("genesis read error: %v", err)
		}
	}
	genesis, err := emtUtils.ParseBlobGenesis(genesisBlob)
	if err != nil {
		ethUtils.Fatalf("genesisJSON err: %v", err)
	}

	genesisFileName := filepath.Join(lightchainDataDir, "genesis.json")
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
	lsCfgPath := filepath.Join(lightchainDataDir, lsConfig.CONFIG_FOLDER, lsConfig.CONFIG_NAME)
	err = ioutil.WriteFile(lsCfgPath, lsConfig.ReadDefaultLsConfigBlob(), 0666)
	if err != nil {
		ethUtils.Fatalf("LsConfig err: %v", err)
	} else {
		ethLog.Info(fmt.Sprintf("successfully copied LS config into: %s", lsCfgPath))
	}
	
	return nil
}

func ResetCmd(ctx *cli.Context) error {
	dbDir := filepath.Join(emtUtils.MakeDataDir(ctx), "lightchain")
	if err := os.RemoveAll(dbDir); err != nil {
		ethLog.Debug("Could not reset lightchain. Failed to remove %+v", dbDir)
		return err
	}

	ethLog.Info("Successfully removed all data", "dir", dbDir)
	return nil
}

func init() {
	app.Action = LightchainNodeCmd // Fallback command
	app.HideVersion = true
	app.Before = BeforeCmd
	app.Commands = []cli.Command{
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Print lightchain and go-ethereum version in usage",
			Action:  VersionCmd,
		},
		{
			Action:      InitCmd,
			Name:        "init",
			Usage:       "init genesis.json",
			Description: "Initialize the files",
		},
		{
			Action: ResetCmd,
			Name:   "unsafe_reset_all",
			Usage:  "(unsafe) Remove lightchain database",
		},
		{
			Action: LightchainNodeCmd,
			Name:   "node",
			Usage:  "Running Lightchain app",
		},
		{
			Action: kvStoreCmd,
			Name:   "kvstore",
			Usage:  "Running abci-kvstore app",
		},
	}
	app.Flags = append(app.Flags, emtUtils.NodeFlags...)
	app.Flags = append(app.Flags, emtUtils.RpcFlags...)
	app.Flags = append(app.Flags, emtUtils.LightchainFlags...)
}

func buildLsConfigPath(ctx *cli.Context) string {
	return filepath.Join(ctx.GlobalString(ethUtils.DataDirFlag.Name), lsConfig.CONFIG_FOLDER, lsConfig.CONFIG_NAME)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
