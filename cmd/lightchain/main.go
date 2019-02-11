package main

import (
	"fmt"
	"os"
	"gopkg.in/urfave/cli.v1"
	"github.com/spf13/cobra"
	"github.com/mitchellh/go-homedir"
	"path"
	"github.com/lightstreams-network/lightchain/log"

	ethLog "github.com/ethereum/go-ethereum/log"
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
)

var logger = log.NewLogger().With("engine", "lightchain")

var (
	defaultHomeDir, _ = homedir.Dir()

	DataDirFlag = ethUtils.DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: ethUtils.DirectoryString{path.Join(defaultHomeDir, "lightchain")},
	}

	LogLvlFlag = cli.StringFlag{
		Name:  "lvl",
		Usage: "Level of logging",
		Value: ethLog.LvlInfo.String(),
	}
)

func main() {
	lightchainCmd := LightchainCmd()

	if err := lightchainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// LightchainCmd is the main Lightstreams PoA blockchain node.
func LightchainCmd() *cobra.Command {
	var lightchainCmd = &cobra.Command{
		Use:   "lightchain",
		Short: "Lightstreams PoA blockchain node.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	lightchainCmd.AddCommand(versionCmd)
	lightchainCmd.AddCommand(docsCmd())
	lightchainCmd.AddCommand(initCmd())
	lightchainCmd.AddCommand(runCmd())
	return lightchainCmd
}

func addDefaultFlags(cmd *cobra.Command) {
	cmd.Flags().String(DataDirFlag.Name, DataDirFlag.Value.Value, DataDirFlag.Usage)
	cmd.Flags().String(LogLvlFlag.Name, LogLvlFlag.Value, LogLvlFlag.Usage)

	cmd.MarkFlagRequired(DataDirFlag.Name)	
}
