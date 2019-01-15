package main

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/spf13/cobra"
	"os"
	"fmt"

	ethLog "github.com/ethereum/go-ethereum/log"

	"github.com/lightstreams-network/lightchain/node"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
	"path/filepath"
	"github.com/lightstreams-network/lightchain/log"
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
	initCmd.Flags().Bool(SiriusNetFlag.Name, true, SiriusNetFlag.Usage)
	return initCmd
}

func initCmdRun(cmd *cobra.Command, args []string) {
	var network = consensus.SiriusNetwork
	lvlStr, _ := cmd.Flags().GetString(LogLvlFlag.Name)
	if lvl, err := ethLog.LvlFromString(lvlStr); err == nil {
		log.SetupLogger(lvl)
	}

	dataDir, _ := cmd.Flags().GetString(DataDirFlag.Name)
	useStandAloneNet, _ := cmd.Flags().GetBool(StandAloneNetFlag.Name)
	
	if useStandAloneNet {
		network = ""
	}
	
	consensusCfg := consensus.NewConfig(
		filepath.Join(dataDir, consensus.DataDirName),
		TendermintRpcListenPort,
		TendermintProxyListenPort,
		TendermintP2PListenPort,
		TendermintProxyProtocol,
	)
	
	dbDataDir := filepath.Join(dataDir, database.DataDirPath)
	ctx := newNodeClientCtx(dbDataDir, cmd)

	dbCfg, err := database.NewConfig(dbDataDir, ctx)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	nodeCfg := node.NewConfig(dataDir, consensusCfg, dbCfg)
	if err := node.InitNode(nodeCfg, network); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	
	logger.Info(fmt.Sprintf("Lightchain node successfully initialized into '%s'!", dataDir))
	os.Exit(0)
}
