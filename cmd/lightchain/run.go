package main

import (
	"github.com/spf13/cobra"
	"path/filepath"
	"fmt"
	"os"
	"gopkg.in/urfave/cli.v1"
	
	ethLog "github.com/ethereum/go-ethereum/log"
	
	"github.com/lightstreams-network/lightchain/node"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/tendermint/tendermint/libs/common"
)

const (
	TendermintP2PListenPort   = uint(26656)
	TendermintRpcListenPort   = uint(26657)
	TendermintProxyListenPort = uint(26658)
	TendermintProxyProtocol   = "rpc"
)

var (
	ConsensusRpcListenPortFlag = cli.UintFlag{
		Name:  "tmt_rpc_port",
		Value: TendermintRpcListenPort,
		Usage: "Tendermint RPC port used to receive incoming messages from Lightchain",
	}
	ConsensusP2PListenPortFlag = cli.UintFlag{
		Name:  "tmt_p2p_port",
		Value: TendermintP2PListenPort,
		Usage: "Tendermint port used to achieve exchange messages across nodes",
	}
	ConsensusProxyListenPortFlag = cli.UintFlag{
		Name:  "tmt_proxy_port",
		Value: TendermintProxyListenPort,
		Usage: "Lightchain RPC port used to receive incoming messages from Tendermint",
	}
	// ABCIProtocolFlag defines whether GRPC or SOCKET should be used for the ABCI connections
	ConsensusProxyProtocolFlag = cli.StringFlag{
		Name:  "abci_protocol",
		Value: "socket",
		Usage: "socket | grpc",
	}
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
					logger.Error(fmt.Errorf("error stopping lightchain node. %v", err).Error())
				}
				os.Exit(1)
			})

			os.Exit(0)
		},
	}

	addDefaultFlags(runCmd)
	addConsensusFlags(runCmd)
	addEthNodeFlags(runCmd)

	return runCmd
}

func addConsensusFlags(cmd *cobra.Command) {
	// Consensus Flags
	cmd.Flags().Uint(ConsensusRpcListenPortFlag.GetName(), ConsensusRpcListenPortFlag.Value, ConsensusRpcListenPortFlag.Usage)
	cmd.Flags().Uint(ConsensusP2PListenPortFlag.GetName(), ConsensusP2PListenPortFlag.Value, ConsensusP2PListenPortFlag.Usage)
	cmd.Flags().Uint(ConsensusProxyListenPortFlag.GetName(), ConsensusProxyListenPortFlag.Value, ConsensusProxyListenPortFlag.Usage)
	cmd.Flags().String(ConsensusProxyProtocolFlag.GetName(), ConsensusProxyProtocolFlag.Value, ConsensusProxyProtocolFlag.Usage)
}

