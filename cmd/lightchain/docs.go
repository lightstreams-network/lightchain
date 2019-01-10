package main

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/spf13/cobra/doc"
	"os"
	"path"
)

const docsCmdDir = "docs/cmd/auto_generated"

func docsCmd() *cobra.Command {
	var docsCmd = &cobra.Command{
		Use:   "docs",
		Short: fmt.Sprintf("Generates lightchain cmd usage docs based on code into the: '%s'.", docsCmdDir),
		Run: func(cmd *cobra.Command, args []string) {
			docsDirPath := getDocsDirPath()

			err := doc.GenMarkdownTree(LightchainCmd(), docsDirPath)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			logger.Info(fmt.Sprintf("Documentation for `lightchain` successfully generated into the '%s'!", docsDirPath))
		},
	}

	return docsCmd
}

func getDocsDirPath() string {
	dir, err := os.Getwd()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	return path.Join(dir, docsCmdDir)
}