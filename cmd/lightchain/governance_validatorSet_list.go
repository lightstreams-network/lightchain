package main

import (
	"os"
	"github.com/spf13/cobra"
	
	"github.com/lightstreams-network/lightchain/node"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/ethereum/go-ethereum/common"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/lightstreams-network/lightchain/governance"
	"fmt"
)


func governanceValidatorSetListCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "validatorset-list",
		Short: "Launch a lightchain node and add a new validator to ValidatorSet contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.SetupLogger(ethLog.LvlInfo)
			logger.Info("Simulating Lightchain node activity to verify state and collect stats...")
			
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

			validatorAddress, err := cmd.Flags().GetString(ValidatorAddressFlag.Name)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			
			pubKeys, err := listValidatorSet(nodeCfg, common.HexToAddress(validatorAddress))
			if err != nil {
				logger.Error(err.Error())
				n.Stop()
				os.Exit(1)
			}
			
			logger.Info(fmt.Sprintf("Active validator pub keys: %+v", pubKeys))
		},
	}
	
	addRunCmdFlags(cmd)
	addGovernanceCmdFlags(cmd)
	cmd.Flags().String(ValidatorPubKeyFlag.GetName(), "", OwnerAccountFlag.Usage)
	cmd.Flags().String(ValidatorAddressFlag.GetName(), "", ValidatorAddressFlag.Usage)
	
	return cmd
}


func listValidatorSet(nodeCfg node.Config, validatorAddr common.Address) ([]string, error) {

	instance := governance.NewValidatorSet(nodeCfg.GovernanceCfg().ContractAddress(), nodeCfg.DbCfg().GethIpcPath())
	return instance.FetchPubKeySet(validatorAddr)
}