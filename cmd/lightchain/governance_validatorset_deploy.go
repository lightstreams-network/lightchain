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
	"time"
	"github.com/lightstreams-network/lightchain/database"
)


func governanceValidatorSetDeployCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "validatorset-deploy",
		Short: "Launch a lightchain node and deploy governance smart contract",
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
			
			
			nodeCfg, err := loadRunCmdConfig(cmd)
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

			validatorSetContractAddress, err := deployValidatorSetContract(nodeCfg, common.HexToAddress(owner), password)
			if err != nil {
				logger.Error(err.Error())
				n.Stop()
				os.Exit(1)
			}
			
			logger.Info("Wait few seconds for block to persist...")
			time.Sleep(time.Second * 2)

			if nodeCfg.TracerCfg().ShouldTrace {
				logger.Info("Running tracer assertion over deployed validatorSet smart contract...")
				assertPostDeployState(nodeCfg, validatorSetContractAddress, common.HexToAddress(owner))
			}

			n.Stop()
			fmt.Printf("\n\nSmart contract was succesfully deployed at %s . \n\n", validatorSetContractAddress.String())
		},
	}
	
	addRunCmdFlags(cmd)
	addGovernanceCmdFlags(cmd)
	
	return cmd
}

func deployValidatorSetContract(nodeCfg node.Config, owner common.Address, password string) (common.Address, error) {
	logger.Info("Deploying ValidatorSet contract...")

	txAuth, err := authy.FindInKeystoreDir(nodeCfg.DbCfg().KeystoreDir(), owner, password)
	if err != nil {
		return common.Address{}, err
	}

	address, err := governance.DeployContract(txAuth, nodeCfg.DbCfg().GethIpcPath())
	if err != nil {
		return common.Address{}, err
	}
	
	return address, nil
}


func assertPostDeployState(nodeCfg node.Config, contractAddress common.Address, owner common.Address) {
	tracer, err := database.NewTracer(nodeCfg.TracerCfg(), nodeCfg.DbCfg().ChainDbDir(), nodeCfg.DbCfg().GethIpcPath())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	tracer.AssertPersistedValidatorSetContract(contractAddress, owner)
}