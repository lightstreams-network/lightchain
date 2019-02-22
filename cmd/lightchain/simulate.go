package main

import (
	"os"
	"time"
	"fmt"
	"math/big"
	"github.com/spf13/cobra"
	"github.com/lightstreams-network/lightchain/node"
	"github.com/lightstreams-network/lightchain/wallety"
	"github.com/lightstreams-network/lightchain/txclient"
	"github.com/lightstreams-network/lightchain/authy"
	"github.com/lightstreams-network/lightchain/log"
	ethLog "github.com/ethereum/go-ethereum/log"
)

const standaloneGenesisAcc = "0xc916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e"
const standaloneGenesisPwd = "WelcomeToSirius"
const standaloneNonExistingAddr = "0xc111111111111111111111111111111111111111"
const simulateTxAfterSeconds = 10

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

			// Give the node few secs to boot and then fire a TX.
			// This should be improved in the future, awaiting the node boot-up via a channel
			logger.Info("Simulating TX to modify state...")
			go simulateTransferTx(nodeCfg)

			err = startNode(nodeCfg, stopSimulatedNode)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			logger.Info("Lightchain simulation finished.")
			os.Exit(0)
		},
	}

	addInitCmdFLags(simulateCmd)
	addRunCmdFlags(simulateCmd)

	return simulateCmd
}

func simulateTransferTx(nodeCfg node.Config) {
	time.Sleep(time.Second*simulateTxAfterSeconds)
	client, err := txclient.Dial(nodeCfg.DbCfg().GethIpcPath())
	if err != nil {
		panic(err)
	}

	auth, err := authy.FindInKeystoreDir(nodeCfg.DbCfg().KeystoreDir(), authy.NewEthAccountFromHex(standaloneGenesisAcc), standaloneGenesisPwd)
	if err != nil {
		panic(err)
	}

	amount := big.NewInt(1e+18)

	err = wallety.Transfer(client, auth, authy.NewEthAccountFromHex(standaloneNonExistingAddr), amount.String())
	if err != nil {
		panic(err)
	}
}

func stopSimulatedNode(n *node.Node) {
	time.Sleep(time.Second*simulateTxAfterSeconds+time.Second*5)
	n.Stop()
}