package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

const Major = "1"
const Minor = "3"
const Fix = "0"
const Verbal = "Fast && Safe"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Describes version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Version: %s.%s.%s %s", Major, Minor, Fix, Verbal))
	},
}
