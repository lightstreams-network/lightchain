package main

import (
	"github.com/spf13/cobra"
	"path/filepath"
	"fmt"
	"os"
	
	ethLog "github.com/ethereum/go-ethereum/log"
	
	"github.com/lightstreams-network/lightchain/node"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/tendermint/tendermint/libs/common"
)

func runCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Launches lightchain node and all of its online services including blockchain (Geth) and the consensus (Tendermint).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			lvlStr, _ := cmd.Flags().GetString(LogLvlFlag.Name)
			if lvl, err := ethLog.LvlFromString(lvlStr); err == nil {
				log.SetupLogger(lvl)
			}
			
			logger.Info("Launching Lightchain node...")

			dataDir, _ := cmd.Flags().GetString(DataDirFlag.GetName())
			rpcListenPort, _ := cmd.Flags().GetUint(ConsensusRpcListenPortFlag.GetName())
			p2pListenPort, _ := cmd.Flags().GetUint(ConsensusP2PListenPortFlag.GetName())
			proxyListenPort, _ := cmd.Flags().GetUint(ConsensusProxyListenPortFlag.GetName())
			proxyProtocol, _ := cmd.Flags().GetString(ConsensusProxyProtocolFlag.GetName())
			databaseDataDir := filepath.Join(dataDir, database.DataDirPath)

			consensusCfg := consensus.NewConfig(
				filepath.Join(dataDir, consensus.DataDirName),
				rpcListenPort,
				p2pListenPort,
				proxyListenPort,
				proxyProtocol,
			)
			
			// Fake cli.context required by Ethereum node
			ctx := newNodeClientCtx(databaseDataDir, cmd)
			dbCfg, err := database.NewConfig(databaseDataDir, ctx)
			if err != nil {
				logger.Error(fmt.Errorf("database node config could not be created: %v", err).Error())
				os.Exit(1)
			}
			
			nodeCfg := node.NewConfig(dataDir, consensusCfg, dbCfg)

			lightChainNode, err := node.NewNode(&nodeCfg)
			if err != nil {
				logger.Error(fmt.Errorf("lightchain node could not be instantiated: %v", err).Error())
				os.Exit(1)
			}
			
			logger.Debug("Starting lightchain node...")
			if err := lightChainNode.Start(); err != nil {
				logger.Error(fmt.Errorf("lightchain node could not be started: %v", err).Error())
				os.Exit(1)
			}
			
			common.TrapSignal(func() {
				if err := lightChainNode.Stop(); err != nil {
					logger.Error(fmt.Errorf("error stopping Tendermint service. %v", err).Error())
				}
				os.Exit(1)
			})

			os.Exit(0)
		},
	}

	addDefaultFlags(runCmd)
	addEthNodeFlags(runCmd)

	return runCmd
}
