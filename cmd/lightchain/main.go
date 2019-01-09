package main

import (
	"fmt"
	"os"
	"gopkg.in/urfave/cli.v1"
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	
	"github.com/lightstreams-network/lightchain/utils"
	"github.com/lightstreams-network/lightchain/log"
)

var (
	// The app that holds all commands and flags.
	app = ethUtils.NewApp(Version, "the lightchain command line interface")
)

func BeforeCmd(ctx *cli.Context) error {
	logLvl := ctx.GlobalInt(utils.VerbosityFlag.Name)
	if err := log.SetupLogger(logLvl); err != nil {
		return err
	}
	return nil
}

func AfterCmd(ctx *cli.Context) error {
	return nil
}

func init() {
	app.Action = RunCmd // Fallback command
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
			Action: RunCmd,
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
