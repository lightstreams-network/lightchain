package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

const Major = "1"
const Minor = "2"
const Fix = "0"
const Verbal = "Rock && Roll"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Describes version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Version: %s.%s.%s %s", Major, Minor, Fix, Verbal))
	},
}
