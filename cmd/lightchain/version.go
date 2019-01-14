package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Major version component of the current release
const Major = "0"

// Minor version component of the current release
const Minor = "9"

// Fix version component of the current release
const Fix = "0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Describes version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Version: %s.%s.%s-alpha Sirius-Net", Major, Minor, Fix))
	},
}
