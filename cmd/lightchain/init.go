package main

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/spf13/cobra"
	"os"
	"fmt"
	"path/filepath"

	ethLog "github.com/ethereum/go-ethereum/log"

	"github.com/lightstreams-network/lightchain/node"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/lightstreams-network/lightchain/setup"
	"github.com/lightstreams-network/lightchain/prometheus"
)

var (
	StandAloneNetFlag = cli.BoolFlag{
		Name:  "standalone",
		Usage: "Initialize a stand alone node not connected to any network",
	}

	SiriusNetFlag = cli.BoolFlag{
		Name:  "sirius",
		Usage: "Initialize a node connected to Sirius network",
	}
)

func initCmd() *cobra.Command {
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initializes new lightchain node according to the configured flags.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: initCmdRun,
	}

	addDefaultFlags(initCmd)
	initCmd.Flags().Bool(StandAloneNetFlag.Name, false, DataDirFlag.Usage)
	initCmd.Flags().Bool(SiriusNetFlag.Name, false, SiriusNetFlag.Usage)

	return initCmd
}

func initCmdRun(cmd *cobra.Command, args []string) {
	lvlStr, _ := cmd.Flags().GetString(LogLvlFlag.Name)
	if lvl, err := ethLog.LvlFromString(lvlStr); err == nil {
		log.SetupLogger(lvl)
	}

	dataDir, _ := cmd.Flags().GetString(DataDirFlag.Name)
	useStandAloneNet, _ := cmd.Flags().GetBool(StandAloneNetFlag.Name)
	useSiriusNet, _ := cmd.Flags().GetBool(SiriusNetFlag.Name)

	shouldTrace, _ := cmd.Flags().GetBool(TraceFlag.Name)
	traceLogFilePath, _ := cmd.Flags().GetString(TraceLogFlag.Name)

	if shouldTrace {
		logger.Info("|--------")
		logger.Info("| Danger: Tracing enabled is not recommended in production!")
		logger.Info(fmt.Sprintf("| Tracing output is configured to be persisted at %v", traceLogFilePath))
		logger.Info("|--------")
	}

	var network setup.Network
	if useStandAloneNet && useSiriusNet {
		logger.Error(fmt.Errorf("multiple networks selected: %s, %s", setup.SiriusNetwork, setup.StandaloneNetwork).Error())
		os.Exit(1)
	} else if useStandAloneNet {
		network = setup.StandaloneNetwork
	} else if useSiriusNet {
		network = setup.SiriusNetwork
	} else {
		network = setup.SiriusNetwork
	}


	consensusCfg := consensus.NewConfig(
		filepath.Join(dataDir, consensus.DataDirName),
		TendermintRpcListenPort,
		TendermintProxyListenPort,
		TendermintP2PListenPort,
		TendermintProxyProtocol,
	)

	dbDataDir := filepath.Join(dataDir, database.DataDirPath)
	dbCfg, err := database.NewConfig(dbDataDir, shouldTrace, traceLogFilePath, newNodeClientCtx(dbDataDir, cmd))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	
	prometheusCfg := prometheus.NewConfig(false,
		prometheus.DefaultPrometheusAddr,
		prometheus.DefaultPrometheusNamespace,
		dbCfg.GethIpcPath())

	nodeCfg := node.NewConfig(dataDir, consensusCfg, dbCfg, prometheusCfg)
	if err := node.Init(nodeCfg, network); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Lightchain node successfully initialized into '%s'!", dataDir))
	os.Exit(0)
}
