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
)

var (
	OwnerAccountFlag = cli.UintFlag{
		Name:  "owner",
		Usage: "ValidatorSet smart contract owner account",
	}
	OwnerAccountPasswordFlag = cli.UintFlag{
		Name:  "password",
		Usage: "Passphrase to decrypt owner account",
	}
)


func governanceCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "governance",
		Short: "Executes `run` commands and deploy governance smart contract using genesis validators",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.SetupLogger(ethLog.LvlDebug)
			logger.Info("Simulating Lightchain node activity to verify state and collect stats...")
			
			owner, _ := cmd.Flags().GetString(OwnerAccountFlag.Name)
			if owner == "" {
				logger.Error(fmt.Sprintf("Missing value for argument %v", OwnerAccountFlag.Name))
			}
			
			password, _ := cmd.Flags().GetString(OwnerAccountPasswordFlag.Name)
			if password == "" {
				password, _ = fs.PromptPassword(fmt.Sprintf("Enter password to decrypt the account %s: ", owner))
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
			defer n.Stop()

			scAddress, err := deployValidatorSetContract(nodeCfg, common.HexToAddress(owner), password)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			logger.Info(fmt.Sprintf("Smart contract was deployed successfully: %s.", scAddress.String()))
			os.Exit(0)
		},
	}
	
	addRunCmdFlags(cmd)
	cmd.Flags().String(OwnerAccountFlag.GetName(), "", OwnerAccountFlag.Usage)
	cmd.Flags().String(OwnerAccountPasswordFlag.GetName(), "", OwnerAccountPasswordFlag.Usage)
	
	return cmd
}

func deployValidatorSetContract(nodeCfg node.Config, owner common.Address, password string) (common.Address, error) {
	logger.Info("Deploying ValidatorSet contract...")

	txAuth, err := authy.FindInKeystoreDir(nodeCfg.DbCfg().KeystoreDir(), owner, password)
	if err != nil {
		return common.Address{}, nil
	}

	return governance.DeployContract(txAuth, nodeCfg.DbCfg().GethIpcPath())
}