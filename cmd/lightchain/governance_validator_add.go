package main

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/urfave/cli.v1"
	
	"github.com/lightstreams-network/lightchain/node"
	"github.com/lightstreams-network/lightchain/authy"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/ethereum/go-ethereum/common"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/lightstreams-network/lightchain/fs"
	"github.com/lightstreams-network/lightchain/governance"
	"time"
)

var (
	ValidatorKeyFlag = cli.StringFlag{
		Name:  "pubkey",
		Usage: "PubKey of validator to append in the ValidatorSet contract",
	}
	ValidatorAddressFlag = cli.StringFlag{
		Name:  "address",
		Usage: "Validator ethereum address to link to added validator",
	}
)


func governanceValidatorAddCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "validator-add",
		Short: "Launch a lightchain node and add a new validator to ValidatorSet contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.SetupLogger(ethLog.LvlInfo)
			logger.Info("Simulating Lightchain node activity to verify state and collect stats...")
			
			owner, _ := cmd.Flags().GetString(OwnerAccountFlag.Name)
			if owner == "" {
				logger.Error(fmt.Sprintf("Missing value for argument %v", OwnerAccountFlag.Name))
				os.Exit(1)
			}
			
			password, _ := cmd.Flags().GetString(OwnerAccountPasswordFlag.Name)
			if password == "" {
				var err error
				password, err = fs.PromptPassword(fmt.Sprintf("Enter password to decrypt the account %s: ", owner))
				if err != nil {
					logger.Error(err.Error())
					os.Exit(1)
				}
			}
			
			validatorPubKey, _ := cmd.Flags().GetString(ValidatorKeyFlag.Name)
			if owner == "" {
				logger.Error(fmt.Sprintf("Missing value for argument %v", ValidatorKeyFlag.Name))
				os.Exit(1)
			}
			
			validatorAddress, _ := cmd.Flags().GetString(ValidatorAddressFlag.Name)
			if owner == "" {
				logger.Error(fmt.Sprintf("Missing value for argument %v", ValidatorAddressFlag.Name))
				os.Exit(1)
			}
			
			
			nodeCfg, err := newRunCmdConfig(cmd)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			n, err := node.NewNode(&nodeCfg)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			if err := n.Start(); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			err = addValidator(nodeCfg, common.HexToAddress(owner), password, validatorPubKey, common.HexToAddress(validatorAddress))
			if err != nil {
				logger.Error(err.Error())
				n.Stop()
				os.Exit(1)
			}
			
			logger.Info("Wait few seconds for block to persist...")
			time.Sleep(time.Second * 2)
			n.Stop()

			fmt.Printf("\n\nValidator %s was linked to %s successfully .\n\n", validatorPubKey, validatorAddress)
		},
	}
	
	addRunCmdFlags(cmd)
	addGovernanceCmdFlags(cmd)
	cmd.Flags().String(ValidatorKeyFlag.GetName(), "", OwnerAccountFlag.Usage)
	cmd.Flags().String(ValidatorAddressFlag.GetName(), "", ValidatorAddressFlag.Usage)
	
	return cmd
}


func addValidator(nodeCfg node.Config, owner common.Address, password string, pubKey string, validatorAddr common.Address) error {
	txAuth, err := authy.FindInKeystoreDir(nodeCfg.DbCfg().KeystoreDir(), owner, password)
	if err != nil {
		return err
	}

	instance := governance.NewValidatorSet(nodeCfg.GovernanceCfg().ContractAddress(), nodeCfg.DbCfg().GethIpcPath())
	return instance.AddValidator(txAuth, pubKey, validatorAddr)
}