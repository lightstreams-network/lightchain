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
	"github.com/lightstreams-network/lightchain/prometheus"
	"github.com/lightstreams-network/lightchain/tracer"
	"github.com/lightstreams-network/lightchain/governance"
	"path"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/lightstreams-network/lightchain/network"
)

const (
	TendermintP2PListenPort = uint(26656)
	TendermintRpcListenPort = uint(26657)
	TendermintProxyAppName  = "lightchain"
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

			nodeCfg, err := newRunCmdConfig(cmd)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			
			n, err := node.NewNode(&nodeCfg)
			if err != nil {
				logger.Error(fmt.Errorf("lightchain node could not be instantiated: %v", err).Error())
				os.Exit(1)
			}

			common.TrapSignal(logger, func() {
				if err := n.Stop(); err != nil {
					logger.Error(fmt.Errorf("error stopping lightchain node. %v", err).Error())
					os.Exit(1)
				}

				os.Exit(0)
			})

			logger.Debug("Starting lightchain node...")
			if err := n.Start(); err != nil {
				logger.Error(fmt.Errorf("lightchain node could not be started: %v", err).Error())
				os.Exit(1)
			}

			select {}
		},
	}

	addRunCmdFlags(runCmd)

	return runCmd
}

func newRunCmdConfig(cmd *cobra.Command) (node.Config, error){
	dataDir, _ := cmd.Flags().GetString(DataDirFlag.GetName())
	shouldTrace, _ := cmd.Flags().GetBool(TraceFlag.Name)
	rpcListenPort, _ := cmd.Flags().GetUint(ConsensusRpcListenPortFlag.GetName())
	p2pListenPort, _ := cmd.Flags().GetUint(ConsensusP2PListenPortFlag.GetName())
	proxyAppName, _ := cmd.Flags().GetString(ConsensusProxyAppNameFlag.GetName())
	enablePrometheus, _ := cmd.Flags().GetBool(PrometheusFlag.GetName())
	databaseDataDir := filepath.Join(dataDir, database.DataDirPath)

	consensusCfg, err := consensus.NewConfig(
		filepath.Join(dataDir, consensus.DataDirName),
		rpcListenPort,
		p2pListenPort,
		proxyAppName,
		enablePrometheus,
	)
	if err != nil {
		return node.Config{}, fmt.Errorf("consensus node config could not be created: %v", err)
	}

	// Fake cli.context required by Ethereum node
	ctx := newNodeClientCtx(databaseDataDir, cmd)
	dbCfg, err := database.NewConfig(databaseDataDir, enablePrometheus, ctx)
	if err != nil {
		return node.Config{}, fmt.Errorf("database node config could not be loaded: %v", err)
	}

	prometheusCfg := prometheus.NewConfig(
		enablePrometheus,
		prometheus.DefaultPrometheusAddr,
		prometheus.DefaultPrometheusNamespace,
		dbCfg.GethIpcPath(),
	)

	tracerCfg := tracer.NewConfig(shouldTrace, path.Join(dataDir, "tracer.log"))

	if shouldTrace {
		tracerCfg.PrintWarning(logger)
	}

	governanceCfg, err := governance.Load(dataDir)
	if err != nil {
		logger.Error(fmt.Errorf("governance config could not be loaded: %v", err).Error())
		networkId, err := dbCfg.NetworkId()
		if err != nil {
			return node.Config{}, nil
		}

		ntw, err := network.NewNetworkFromId(networkId)
		if err != nil {
			return node.Config{}, nil
		}
		governanceCfg, err = governance.Init(ntw, dataDir)
		if err != nil {
			return node.Config{}, nil
		}
	}

	nodeCfg := node.NewConfig(dataDir, consensusCfg, dbCfg, prometheusCfg, tracerCfg, governanceCfg)
	return nodeCfg, nil
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
}
