package main

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	
	"github.com/lightstreams-network/lightchain/node"
	"github.com/lightstreams-network/lightchain/authy"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/ethereum/go-ethereum/common"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/lightstreams-network/lightchain/fs"
	"github.com/lightstreams-network/lightchain/governance"
	"github.com/lightstreams-network/lightchain/database"
	"time"
)

func governanceValidatorSetAddCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "validatorset-add",
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
			
			validatorPubKey, _ := cmd.Flags().GetString(ValidatorPubKeyFlag.Name)
			if validatorPubKey == "" {
				logger.Error(fmt.Sprintf("Missing value for argument '%v'", ValidatorPubKeyFlag.Name))
				os.Exit(1)
			}
			
			validatorAddress, _ := cmd.Flags().GetString(ValidatorAddressFlag.Name)
			if validatorAddress == "" {
				logger.Error(fmt.Sprintf("Missing value for argument '%v'", ValidatorAddressFlag.Name))
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
			time.Sleep(time.Second * 5)
			n.Stop()

			if nodeCfg.TracerCfg().ShouldTrace {
				logger.Info("Running tracer assertion over added validator action...")
				assertPostAddValidatorState(nodeCfg, validatorPubKey, common.HexToAddress(validatorAddress))
			}

			fmt.Printf("\n\nValidator %s:%s was added successfully.\n\n", validatorPubKey, validatorAddress)
		},
	}
	
	addRunCmdFlags(cmd)
	addGovernanceCmdFlags(cmd)
	cmd.Flags().String(ValidatorPubKeyFlag.GetName(), "", OwnerAccountFlag.Usage)
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

func assertPostAddValidatorState(nodeCfg node.Config, validatorPubKey string, validatorAddress common.Address) {
	tracer, err := database.NewTracer(nodeCfg.TracerCfg(), nodeCfg.DbCfg().ChainDbDir(), nodeCfg.DbCfg().GethIpcPath())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	tracer.AssertPersistedValidatorSetAddValidator(nodeCfg.ConsensusCfg().TendermintCfg(), validatorPubKey, validatorAddress)
}