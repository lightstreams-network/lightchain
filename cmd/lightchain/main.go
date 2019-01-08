package main

import (
	"fmt"
	"os"
	"gopkg.in/urfave/cli.v1"
	"path/filepath"

	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	ethLog "github.com/ethereum/go-ethereum/log"
	
	tmtServer "github.com/tendermint/tendermint/abci/server"
	tmtLog "github.com/tendermint/tendermint/libs/log"
	tmtCommon "github.com/tendermint/tmlibs/common"

	"github.com/lightstreams-network/lightchain/utils"
	"github.com/lightstreams-network/lightchain/config"
	"github.com/lightstreams-network/lightchain/ethereum"
	"github.com/lightstreams-network/lightchain/version"
	"github.com/lightstreams-network/lightchain/abci/transaction"
	"github.com/lightstreams-network/lightchain/abci"
	"github.com/lightstreams-network/lightchain/tendermint"
)

var (
	// The app that holds all commands and flags.
	app = ethUtils.NewApp(version.Version, "the lightchain command line interface")
)

func BeforeCmd(ctx *cli.Context) error {
	logLvl := ctx.GlobalInt(utils.VerbosityFlag.Name)
	if err := utils.SetupLogger(logLvl); err != nil {
		return err
	}
	return nil
}

func AfterCmd(ctx *cli.Context) error {
	return nil
}

func LightchainNodeCmd(ctx *cli.Context) {
	abciNode := abci.CreateNode(ctx)
	abci.StartNode(ctx, abciNode)

	// Fetch the registered service of this type
	var ethBackend *ethereum.Backend
	if err := abciNode.Service(&ethBackend); err != nil {
		ethUtils.Fatalf("ethereum ethBackend service not running: %v", err)
	}

	// In-proc RPC connection so ABCI Query can be forwarded over the ethereum rpc
	rpcClient, err := abciNode.Attach()
	if err != nil {
		ethUtils.Fatalf("Failed to attach to the inproc geth: %v", err)
	}

	txManager, err := transaction.NewManager(config.ConfigPath(ctx))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create the ABCI application - in memory or persisted to disk
	ethApp, err := abci.CreateLightchainApplication(ethBackend, rpcClient, nil, txManager)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ethApp.SetLogger(utils.LightchainLogger().With("module", "lightchain"))

	// Start the app on the ABCI server listener
	//abciAddr := ctx.GlobalString(utils.ABCIAddrFlag.Name)
	abciAddr := fmt.Sprintf("tcp://0.0.0.0:%d", ctx.GlobalInt(utils.ProxyListenPortFlag.Name))
	abciProtocol := ctx.GlobalString(utils.ABCIProtocolFlag.Name)
	abciSrv, err := tmtServer.NewServer(abciAddr, abciProtocol, ethApp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	abciLogger := tmtLog.NewTMLogger(tmtLog.NewSyncWriter(os.Stdout))
	abciSrv.SetLogger(abciLogger.With("module", "abci-server"))

	if err := abciSrv.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	tmtLogger := tmtLog.NewTMLogger(tmtLog.NewSyncWriter(os.Stdout))
	tmtSrv, err := tendermint.CreateNewNode(ctx,  tmtLogger.With("module", "tendermint"))
	if err != nil {
		fmt.Errorf("Failed to create node: %v", err)
	}
	if err := tendermint.StartNode(ctx, tmtSrv); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tmtCommon.TrapSignal(func() {
		if err := abciNode.Stop(); err != nil {
			fmt.Errorf("Error stopping Geth Node", err)
		}
		if err := abciSrv.Stop(); err != nil {
			fmt.Errorf("Error stopping ABCI service", err)
		}
		if err := tmtSrv.Stop(); err != nil {
			fmt.Errorf("Error stopping Tendermint service", err)
		}
	})
}

func VersionCmd(ctx *cli.Context) error {
	fmt.Println("Version: ", version.Version)
	return nil
}

func InitCmd(ctx *cli.Context) error {
	if err := abci.InitNode(ctx); err != nil {
		return err
	}
	
	if err := ethereum.InitNode(ctx); err != nil {
		return err
	}

	if err := tendermint.InitNode(ctx); err != nil {
		return err
	}
	
	return nil
}

func ResetCmd(ctx *cli.Context) error {
	dbDir := filepath.Join(config.MakeHomeDir(ctx), "lightchain")
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
	app.After = AfterCmd
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
	}

	app.Flags = append(app.Flags, NodeFlags...)
	app.Flags = append(app.Flags, RpcFlags...)
	app.Flags = append(app.Flags, LightchainFlags...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
