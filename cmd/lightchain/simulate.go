package main

import (
	"os"
	"fmt"
	"math/big"
	"github.com/spf13/cobra"
	"github.com/lightstreams-network/lightchain/node"
	"github.com/lightstreams-network/lightchain/wallety"
	"github.com/lightstreams-network/lightchain/txclient"
	"github.com/lightstreams-network/lightchain/authy"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethLog "github.com/ethereum/go-ethereum/log"
)

const simulateTxFrom = "0xc916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e"
const simulateTxFromPwd = "WelcomeToSirius"
const standaloneNonExistingAddr = "0xc111111111111111111111111111111111111111"

var simulateTxAmount = big.NewInt(1e+18)

func simulateCmd() *cobra.Command {
	var simulateCmd = &cobra.Command{
		Use:   "simulate",
		Short: "Executes `init` and `run` commands with active tracing and simulates TXs activity to assert crucial components such as Consensus State, DB, Mempool and others. (testing purposes)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.SetupLogger(ethLog.LvlDebug)
			logger.Info("Simulating Lightchain node activity to verify state and collect stats...")

			useStandAloneNet, _ := cmd.Flags().GetBool(StandAloneNetFlag.Name)

			if !useStandAloneNet {
				logger.Error(fmt.Errorf("simulation only possible in --standalone mode").Error())
				os.Exit(1)
			}

			nodeCfg, _, err := newNodeCfgFromCmd(cmd)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			n, err := startNode(nodeCfg)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			tx := simulateTransferTx(nodeCfg)

			n.Stop()

			assertPostSimulationState(nodeCfg, tx)

			logger.Info("Lightchain simulation finished.")
			os.Exit(0)
		},
	}

	addInitCmdFLags(simulateCmd)
	addRunCmdFlags(simulateCmd)

	return simulateCmd
}

func simulateTransferTx(nodeCfg node.Config) *types.Transaction {
	logger.Info("Simulating 1 TX to modify state...")

	client, err := txclient.Dial(nodeCfg.DbCfg().GethIpcPath())
	if err != nil {
		panic(err)
	}

	auth, err := authy.FindInKeystoreDir(nodeCfg.DbCfg().KeystoreDir(), authy.NewEthAccountFromHex(simulateTxFrom), simulateTxFromPwd)
	if err != nil {
		panic(err)
	}

	tx, err := wallety.Transfer(client, auth, authy.NewEthAccountFromHex(standaloneNonExistingAddr), simulateTxAmount.String())
	if err != nil {
		panic(err)
	}

	return tx
}

func assertPostSimulationState(nodeCfg node.Config, tx *types.Transaction) {
	tracer, err := database.NewTracer(nodeCfg.TracerCfg(), nodeCfg.DbCfg().ChainDbDir())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	tracer.AssertPostTxSimulationState(common.HexToAddress(simulateTxFrom), tx)
}