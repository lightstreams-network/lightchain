package main

import (
	"github.com/spf13/cobra"
	"path/filepath"
	"fmt"
	"os"
	
	tmtCommon "github.com/tendermint/tmlibs/common"
	ethLog "github.com/ethereum/go-ethereum/log"
	
	"github.com/lightstreams-network/lightchain/node"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/log"
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
			
			logger.Info("Launching the lightchain node...")
			dataDir, _ := cmd.Flags().GetString(DataDirFlag.GetName())

			rpcListenPort, _ := cmd.Flags().GetUint(ConsensusRpcListenPortFlag.GetName())
			p2pListenPort, _ := cmd.Flags().GetUint(ConsensusP2PListenPortFlag.GetName())
			proxyListenPort, _ := cmd.Flags().GetUint(ConsensusProxyListenPortFlag.GetName())
			proxyProtocol, _ := cmd.Flags().GetString(ConsensusProxyProtocolFlag.GetName())

			consensusCfg := consensus.NewConfig(
				filepath.Join(dataDir, consensus.DataDirName),
				rpcListenPort,
				p2pListenPort,
				proxyListenPort,
				proxyProtocol,
			)
			
			// Fake cli.context required by Ethereum node
			databaseDataDir := filepath.Join(dataDir, database.DataDirPath)
			ctx := NewNodeClientCtx(databaseDataDir, cmd)
			dbCfg, err := database.NewConfig(databaseDataDir, ctx)
			if err != nil {
				logger.Error(fmt.Errorf("database node config could not be created: %v", err).Error())
				os.Exit(1)
			}
			
			nodeCfg := node.NewConfig(dataDir, consensusCfg, dbCfg)
			logger.Debug("Initializing lightchain node...")
			lightChainNode, err := node.NewNode(&nodeCfg) // Former abciNode
			if err != nil {
				logger.Error(fmt.Errorf("lightchain node could not be initialized: %v", err).Error())
				os.Exit(1)
			}
			
			logger.Debug("Starting lightchain node...")
			if err := lightChainNode.Start(); err != nil {
				logger.Error(fmt.Errorf("lightchain node could not be started: %v", err).Error())
				os.Exit(1)
			}
			
			tmtCommon.TrapSignal(func() {
				if err := lightChainNode.Stop(); err != nil {
					logger.Error(fmt.Errorf("error stopping Tendermint service. %v", err).Error())
				}
				os.Exit(1)
			})

			os.Exit(0)
		},
	}

	addDefaultFlags(runCmd)
	addNodeFlags(runCmd)

	return runCmd
}

func addNodeFlags(cmd *cobra.Command) {
	// RPC Flags
	cmd.Flags().Bool(RPCEnabledFlag.GetName(), false, RPCEnabledFlag.Usage)
	cmd.Flags().String(RPCListenAddrFlag.GetName(), RPCListenAddrFlag.Value, RPCListenAddrFlag.Usage)
	cmd.Flags().Int(RPCPortFlag.GetName(), RPCPortFlag.Value, RPCPortFlag.Usage)
	cmd.Flags().String(RPCApiFlag.GetName(), RPCApiFlag.Value, RPCApiFlag.Usage)
	
	// WS Flags
	cmd.Flags().Bool(WSEnabledFlag.GetName(), false, WSEnabledFlag.Usage)
	cmd.Flags().String(WSListenAddrFlag.GetName(), WSListenAddrFlag.Value, WSListenAddrFlag.Usage)
	cmd.Flags().Int(WSPortFlag.GetName(), WSPortFlag.Value, WSPortFlag.Usage)
	
	// Consensus Flags
	cmd.Flags().Uint(ConsensusRpcListenPortFlag.GetName(), ConsensusRpcListenPortFlag.Value, ConsensusRpcListenPortFlag.Usage)
	cmd.Flags().Uint(ConsensusP2PListenPortFlag.GetName(), ConsensusP2PListenPortFlag.Value, ConsensusP2PListenPortFlag.Usage)
	cmd.Flags().Uint(ConsensusProxyListenPortFlag.GetName(), ConsensusProxyListenPortFlag.Value, ConsensusProxyListenPortFlag.Usage)
	cmd.Flags().String(ConsensusProxyProtocolFlag.GetName(), ConsensusProxyProtocolFlag.Value, ConsensusProxyProtocolFlag.Usage)
}
