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
	ValidatorPubKeyFlag = cli.StringFlag{
		Name:  "pubkey",
		Usage: "PubKey of validator to append in the ValidatorSet contract",
	}
	ValidatorAddressFlag = cli.StringFlag{
		Name:  "address",
		Usage: "Validator ethereum address to link to added validator",
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

	cmd.AddCommand(governanceValidatorSetDeployCmd())
	cmd.AddCommand(governanceValidatorSetAddCmd())
	cmd.AddCommand(governanceValidatorSetRemoveCmd())
	cmd.AddCommand(governanceValidatorSetListCmd())

	return cmd
}

func addGovernanceCmdFlags(cmd *cobra.Command) {
	cmd.Flags().String(OwnerAccountFlag.GetName(), "", OwnerAccountFlag.Usage)
	cmd.Flags().String(OwnerAccountPasswordFlag.GetName(), "", OwnerAccountPasswordFlag.Usage)
}