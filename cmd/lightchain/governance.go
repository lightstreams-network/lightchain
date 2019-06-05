package main

import (
	"github.com/spf13/cobra"
	"gopkg.in/urfave/cli.v1"
)

var (
	OwnerAccountFlag = cli.StringFlag{
		Name:  "owner",
		Usage: "ValidatorSet contract owner account",
	}
	OwnerAccountPasswordFlag = cli.StringFlag{
		Name:  "password",
		Usage: "Passphrase to decrypt owner account",
	}
)

func governanceCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "governance",
		Short: "Manage governance contact content.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return errorIncorrectCmdUsage()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmd.AddCommand(governanceDeployCmd())
	cmd.AddCommand(governanceValidatorAddCmd())

	return cmd
}

func addGovernanceCmdFlags(cmd *cobra.Command) {
	cmd.Flags().String(OwnerAccountFlag.GetName(), "", OwnerAccountFlag.Usage)
	cmd.Flags().String(OwnerAccountPasswordFlag.GetName(), "", OwnerAccountPasswordFlag.Usage)
}