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
	"github.com/lightstreams-network/lightchain/prometheus"
	"github.com/lightstreams-network/lightchain/tracer"
)

const (
	TendermintP2PListenPort   = uint(26656)
	TendermintRpcListenPort   = uint(26657)
	TendermintProxyListenPort = uint(26658)
	TendermintProxyProtocol   = "socket"
	TendermintProxyAppName    = "lightchain"
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
		Value: TendermintProxyProtocol,
		Usage: "socket | grpc",
	}
	ConsensusProxyAppNameFlag = cli.StringFlag{
		Name:  "abci_name",
		Value: TendermintProxyAppName,
		Usage: "socket | grpc",
	}
	PrometheusFlag = cli.BoolFlag{
		Name:  "prometheus",
		Usage: "Enable prometheus metrics exporter",
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
			shouldTrace, _ := cmd.Flags().GetBool(TraceFlag.Name)
			traceLogFilePath, _ := cmd.Flags().GetString(TraceLogFlag.Name)
			rpcListenPort, _ := cmd.Flags().GetUint(ConsensusRpcListenPortFlag.GetName())
			p2pListenPort, _ := cmd.Flags().GetUint(ConsensusP2PListenPortFlag.GetName())
			proxyAppName, _ := cmd.Flags().GetString(ConsensusProxyAppNameFlag.GetName())
			enablePrometheus, _ := cmd.Flags().GetBool(PrometheusFlag.GetName())
			databaseDataDir := filepath.Join(dataDir, database.DataDirPath)

			consensusCfg := consensus.NewConfig(
				filepath.Join(dataDir, consensus.DataDirName),
				rpcListenPort,
				p2pListenPort,
				proxyAppName,
				enablePrometheus,
			)

			// Fake cli.context required by Ethereum node
			ctx := newNodeClientCtx(databaseDataDir, cmd)
			dbCfg, err := database.NewConfig(databaseDataDir, enablePrometheus, ctx)
			if err != nil {
				logger.Error(fmt.Errorf("database node config could not be created: %v", err).Error())
				os.Exit(1)
			}

			prometheusCfg := prometheus.NewConfig(
				enablePrometheus,
				prometheus.DefaultPrometheusAddr,
				prometheus.DefaultPrometheusNamespace,
				dbCfg.GethIpcPath(),
			)

			tracerCfg := tracer.NewConfig(shouldTrace, traceLogFilePath)
			nodeCfg := node.NewConfig(dataDir, consensusCfg, dbCfg, prometheusCfg, tracerCfg)

			n, err := node.NewNode(&nodeCfg)
			if err != nil {
				logger.Error(fmt.Errorf("lightchain node could not be instantiated: %v", err).Error())
				os.Exit(1)
			}

			go func() {
				logger.Debug("Starting lightchain node...")
				if err := n.Start(); err != nil {
					logger.Error(fmt.Errorf("lightchain node could not be started: %v", err).Error())
					os.Exit(1)
				}
			}()

			common.TrapSignal(func() {
				if err := n.Stop(); err != nil {
					logger.Error(fmt.Errorf("error stopping lightchain node. %v", err).Error())
					os.Exit(1)
				}
			})

			os.Exit(0)
		},
	}

	addRunCmdFlags(runCmd)

	return runCmd
}

func addRunCmdFlags(cmd *cobra.Command) {
	addDefaultFlags(cmd)
	addConsensusFlags(cmd)
	addEthNodeFlags(cmd)
	cmd.Flags().Bool(PrometheusFlag.GetName(), false, PrometheusFlag.Usage)
}

func addConsensusFlags(cmd *cobra.Command) {
	cmd.Flags().Uint(ConsensusRpcListenPortFlag.GetName(), ConsensusRpcListenPortFlag.Value, ConsensusRpcListenPortFlag.Usage)
	cmd.Flags().Uint(ConsensusP2PListenPortFlag.GetName(), ConsensusP2PListenPortFlag.Value, ConsensusP2PListenPortFlag.Usage)
	cmd.Flags().Uint(ConsensusProxyListenPortFlag.GetName(), ConsensusProxyListenPortFlag.Value, ConsensusProxyListenPortFlag.Usage)
	cmd.Flags().String(ConsensusProxyProtocolFlag.GetName(), ConsensusProxyProtocolFlag.Value, ConsensusProxyProtocolFlag.Usage)
}