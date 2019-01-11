package main

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/lightstreams-network/lightchain/node"
	"fmt"
	"os"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/database"
	
	tmtCommon "github.com/tendermint/tmlibs/common"
	"github.com/spf13/cobra"
	"path/filepath"
)

func runCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Launches lightchain node and all of its online services including blockchain (Geth) and the consensus (Tendermint).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
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
			ctx := newNodeClientCtx(cmd)
			dbCfg, err := database.NewConfigNode(filepath.Join(dataDir, database.DataDirPath), ctx)
			if err != nil {
				logger.Error(fmt.Errorf("database node config could not be created: %v", err).Error())
				os.Exit(1)
			}
			
			nodeCfg := node.NewConfig(dataDir, consensusCfg, dbCfg)
			lightChainNode, err := node.NewNode(&nodeCfg) // Former abciNode
			if err != nil {
				logger.Error(fmt.Errorf("lightchain node could not be initialized: %v", err).Error())
				os.Exit(1)
			}
			
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

func newNodeClientCtx(cmd *cobra.Command) *cli.Context {
	ctx := cli.NewContext(nil, nil, nil)
	rpcEnabledFlagValue, _ := cmd.Flags().GetBool(RPCEnabledFlag.GetName())

	if rpcEnabledFlagValue {
		ctx.Set(RPCEnabledFlag.GetName(), "true")
	} else {
		ctx.Set(RPCEnabledFlag.GetName(), "false")
	}
	
	rpcListenAddrValue, _ := cmd.Flags().GetString(RPCListenAddrFlag.GetName())
	ctx.Set(RPCListenAddrFlag.GetName(), rpcListenAddrValue)
	
	rpcPortFlagValue, _ := cmd.Flags().GetString(RPCPortFlag.GetName())
	ctx.Set(RPCPortFlag.GetName(), rpcPortFlagValue)
	
	rpcApiFlag, _ := cmd.Flags().GetString(RPCApiFlag.GetName())
	ctx.Set(RPCApiFlag.GetName(), rpcApiFlag)
	
	wsEnabledFlagValue, _ := cmd.Flags().GetBool(WSEnabledFlag.GetName())
	if wsEnabledFlagValue {
		ctx.Set(WSEnabledFlag.GetName(), "true")
	} else {
		ctx.Set(WSEnabledFlag.GetName(), "false")
	}
	
	wsListenAddrFlag, _ := cmd.Flags().GetString(WSListenAddrFlag.GetName())
	ctx.Set(WSListenAddrFlag.GetName(), wsListenAddrFlag)
	
	wsPortFlag, _ := cmd.Flags().GetString(WSPortFlag.GetName())
	ctx.Set(WSPortFlag.GetName(), wsPortFlag)
	
	return ctx
}
