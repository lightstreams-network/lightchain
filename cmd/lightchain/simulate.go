package main

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/lightstreams-network/lightchain/node"
	"time"
)

func simulateCmd() *cobra.Command {
	var simulateCmd = &cobra.Command{
		Use:   "simulate",
		Short: "Executes `init` and `run` commands with active tracing and simulates TXs activity to assert crucial components such as Consensus State, DB, Mempool and others. (testing purposes)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Simulating Lightchain node activity to verify state and collect stats...")

			nodeCfg, ntw, err := newNodeCfgFromCmd(cmd)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			if err := node.Init(nodeCfg, ntw); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

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

func stopSimulatedNode(n *node.Node) {
	time.Sleep(time.Second*15)
	n.Stop()
}