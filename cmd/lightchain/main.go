package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/lightstreams-network/lightchain/log"
)

var logger = log.NewLogger()

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