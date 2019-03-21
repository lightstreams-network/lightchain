package main

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/lightstreams-network/lightchain/distribution"
	"github.com/lightstreams-network/lightchain/log"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/lightstreams-network/lightchain/authy"
	"github.com/ethereum/go-ethereum/common"
)

const (
	distributeFlagCsv = "csv"
	distributeContractHexAddr = "contract_hex_addr"
	distributeGethIpcPath = "geth_ipc_path"
)

const (
	distributeDeployFromKeystoreFilePath = "from_keystore_file_path"
)

func distributeCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "distribute",
		Short: "Loads tokens distribution information from CSV and distributes them through releasable smart contract with vesting.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.SetupLogger(ethLog.LvlInfo)
			csvFilePath, _ := cmd.Flags().GetString(distributeFlagCsv)
			contractHexAddr, _ := cmd.Flags().GetString(distributeContractHexAddr)

			contractAddr := common.HexToAddress(contractHexAddr)

			logger.Info("Distributing tokens from CSV...", "csv", csvFilePath)

			distributionsCount, err := distribution.DistributeFromCsv(csvFilePath, contractAddr)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			logger.Info("Tokens successfully distributed!", "distributions", distributionsCount)
			os.Exit(0)
		},
	}

	cmd.Flags().String(distributeFlagCsv, "", "Absolute file path to a CSV source")
	cmd.MarkFlagRequired(distributeFlagCsv)

	cmd.Flags().String(distributeContractHexAddr, "", "The distribution smart contract hex address where tokens will be locked and ready for token holders to release them to their own personal software wallets")
	cmd.MarkFlagRequired(distributeContractHexAddr)

	cmd.AddCommand(distributeDeployCmd())

	return cmd
}

func distributeDeployCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploys releasable smart contract with vesting distribution.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.SetupLogger(ethLog.LvlInfo)
			keystoreFilePath, _ := cmd.Flags().GetString(distributeDeployFromKeystoreFilePath)
			gethIpcPath, _ := cmd.Flags().GetString(distributeGethIpcPath)

			logger.Info("Deploying distribution contract...")

			logger.Info(fmt.Sprintf("Enter password to decrypt '%s' for contract deployment purposes:", keystoreFilePath))
			password, err := promptPassword()
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			auth, err := authy.NewFromKeystoreFile(keystoreFilePath, password)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			contractAddr, err := distribution.Deploy(auth, gethIpcPath)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			logger.Info("Distribution contract successfully deployed!", "contract_addr", contractAddr)
			os.Exit(0)
		},
	}

	cmd.Flags().String(distributeDeployFromKeystoreFilePath, "", "Path to the distribution contract owner keystore file to authorize deploy")
	cmd.MarkFlagRequired(distributeDeployFromKeystoreFilePath)

	cmd.Flags().String(distributeGethIpcPath, "", "Absolute path to node's geth IPC file")
	cmd.MarkFlagRequired(distributeGethIpcPath)

	return cmd
}