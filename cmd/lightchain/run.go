package main

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/lightstreams-network/lightchain/node"
	"fmt"
	"os"
	"github.com/lightstreams-network/lightchain/utils"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/database"
	
	tmtServer "github.com/tendermint/tendermint/abci/server"
	tmtCommon "github.com/tendermint/tmlibs/common"
	"github.com/lightstreams-network/lightchain/log"
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
			
			
			rpcListenPort, _ := cmd.Flags().GetUint16(ConsensusRpcListenPortFlag.GetName())
			p2pListenPort, _ := cmd.Flags().GetUint16(ConsensusP2PListenPortFlag.GetName())
			proxyListenPort, _ := cmd.Flags().GetUint16(ConsensusProxyListenPortFlag.GetName())
			consensusCfg := consensus.NewConfig(
				filepath.Join(dataDir, consensus.DataDirPath),
				rpcListenPort,
				p2pListenPort,
				proxyListenPort,
			)
			
			// Fake cli.context required by Ethereum node 
			ctx := cli.NewContext(nil, nil, nil)
			dbCfg, err := database.NewConfigNode(
				filepath.Join(dataDir, database.DataDirPath),
				ctx,
			)
			if err != nil {
				logger.Error(fmt.Errorf("database node config could not be created: %v", err).Error())
				os.Exit(1)
			}
			
			nodeCfg := node.NewConfig(dataDir, consensusCfg, dbCfg)
			abciNode, err := node.NewNode(&nodeCfg)
			if err != nil {
				logger.Error(fmt.Errorf("lightchain node could not be created: %v", err).Error())
				os.Exit(1)
			}
			node.StartNode()
			
			
			// @TODO REFACTOR FROM HERE

			// Fetch the registered service of this type
			var ethBackend *database.Database
			if err := abciNode.Service(&ethBackend); err != nil {
				logger.Error(fmt.Errorf("database ethBackend service not running: %v", err).Error())
				os.Exit(1)
			}

			// In-proc RPC connection so ABCI Query can be forwarded over the database rpc
			rpcClient, err := abciNode.Attach()
			if err != nil {
				logger.Error(fmt.Errorf("failed to attach to the inproc geth: %v", err).Error())
				os.Exit(1)
			}

			ethLogger := log.NewLogger()
			ethLogger.With("module", "lightchain")

			// Create the ABCI application - in memory or persisted to disk
			tendermintABCI, err := consensus.NewTendermintABCI(ethBackend, rpcClient, ethLogger)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			// Init the ETH state
			err = tendermintABCI.InitEthState()
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			// Start the app on the ABCI server listener
			abciAddr := fmt.Sprintf("tcp://0.0.0.0:%d", ctx.GlobalInt(utils.TendermintProxyListenPortFlag.Name))
			abciProtocol := ctx.GlobalString(utils.ABCIProtocolFlag.Name)
			abciSrv, err := tmtServer.NewServer(abciAddr, abciProtocol, tendermintABCI)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			abciLogger := log.NewLogger()
			abciSrv.SetLogger(abciLogger.With("module", "node-server"))

			if err := abciSrv.Start(); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			tmtCommon.TrapSignal(func() {
				if err := tmtNode.Stop(); err != nil {
					logger.Error(fmt.Errorf("error stopping Tendermint service. %v", err).Error())
				}
				if err := abciNode.Stop(); err != nil {
					logger.Error(fmt.Errorf("error stopping Geth node. %v", err).Error())
				}
				if err := abciSrv.Stop(); err != nil {
					logger.Error(fmt.Errorf("error stopping ABCI service. %v", err).Error())
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
	cmd.Flags().Int(RPCPortFlag.GetName(), RPCPortFlag.Value, RPCListenAddrFlag.Usage)
	cmd.Flags().String(RPCApiFlag.GetName(), RPCApiFlag.Value, RPCApiFlag.Usage)
	
	// WS Flags
	cmd.Flags().Bool(WSEnabledFlag.GetName(), false, WSEnabledFlag.Usage)
	cmd.Flags().String(WSListenAddrFlag.GetName(), WSListenAddrFlag.Value, WSListenAddrFlag.Usage)
	cmd.Flags().Int(WSPortFlag.GetName(), WSPortFlag.Value, WSPortFlag.Usage)
	
	// Consensus Flags
	cmd.Flags().Uint(ConsensusRpcListenPortFlag.GetName(), ConsensusRpcListenPortFlag.Value, ConsensusRpcListenPortFlag.Usage)
	cmd.Flags().Uint(ConsensusP2PListenPortFlag.GetName(), ConsensusP2PListenPortFlag.Value, ConsensusP2PListenPortFlag.Usage)
	cmd.Flags().Uint(ConsensusProxyListenPortFlag.GetName(), ConsensusProxyListenPortFlag.Value, ConsensusProxyListenPortFlag.Usage)
}
