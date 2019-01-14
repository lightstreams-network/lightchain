package main

import (
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
	return initCmd
}

func initCmdRun(cmd *cobra.Command, args []string) {
	lvlStr, _ := cmd.Flags().GetString(LogLvlFlag.Name)
	if lvl, err := ethLog.LvlFromString(lvlStr); err == nil {
		log.SetupLogger(lvl)
	}
			
	dataDir, _ := cmd.Flags().GetString(DataDirFlag.Name)

	consensusCfg := consensus.NewConfig(
		filepath.Join(dataDir, consensus.DataDirName),
		TendermintRpcListenPort,
		TendermintProxyListenPort,
		TendermintP2PListenPort,
		TendermintProxyProtocol,
	)
	
	dbDataDir := filepath.Join(dataDir, database.DataDirPath)
	ctx := NewNodeClientCtx(dbDataDir, cmd)
	dbCfg, err := database.NewConfig(dbDataDir, ctx)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	nodeCfg := node.NewConfig(dataDir, consensusCfg, dbCfg)
	if err := node.InitNode(nodeCfg); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	
	logger.Info(fmt.Sprintf("Lightchain node successfully initialized into '%s'!", dataDir))
	os.Exit(0)
}
