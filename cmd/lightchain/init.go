package main

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
)

func initCmd() *cobra.Command {
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initializes new lightchain node according to the configured flags.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, _ := cmd.Flags().GetString(flagDataDir)

			//if err := node.InitNode(ctx); err != nil {
			//	logger.Error(err.Error())
			//	os.Exit(1)
			//}
			//
			//if err := database.Init(ctx); err != nil {
			//	logger.Error(err.Error())
			//	os.Exit(1)
			//}
			//
			//if err := consensus.InitNode(ctx); err != nil {
			//	logger.Error(err.Error())
			//	os.Exit(1)
			//}

			logger.Info(fmt.Sprintf("Lightchain node successfully initialized into '%s'!", dataDir))
			os.Exit(0)
		},
	}

	addDefaultFlags(initCmd)

	return initCmd
}