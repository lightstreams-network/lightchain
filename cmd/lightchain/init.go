package main

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"

	"github.com/ethereum/go-ethereum/eth"

	"github.com/lightstreams-network/lightchain/node"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/utils"
	"path/filepath"
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
	dataDir, _ := cmd.Flags().GetString(DataDirFlag.Name)

	consensusCfg := consensus.NewConfig(filepath.Join(dataDir, consensus.DataDirPath),
		uint16(utils.TendermintRpcListenPort),
		uint16(utils.ProxyListenPort), 
		uint16(utils.TendermintP2PListenPort))
	
	dbCfg := database.NewConfig(dataDir, eth.DefaultConfig, database.DefaultEthNodeConfig(), "")

	nodeCfg := node.NewConfig(dataDir, consensusCfg, dbCfg)
	if err := node.InitNode(nodeCfg); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	
	logger.Info(fmt.Sprintf("Lightchain node successfully initialized into '%s'!", dataDir))
	os.Exit(0)
}


//func createDatabaseConfig(dataDir string) (database.Config, error) {
//	ethCfg := eth.DefaultConfig
//	nodeCfg := database.DefaultEthNodeConfig()
//	
//	ethUtils.SetNodeConfig(ctx, &nodeCfg)
//	database.SetNodeDefaultConfig(&nodeCfg, dataDir)
// 	_, err := database.NewNode(&nodeCfg)
// 	if err != nil {
//		return database.Config{}, err
//	}
// 	ethUtils.SetEthConfig(ctx, &stack.Node, &ethCfg)
//	database.SetEthDefaultConfig(&ethCfg)
// 	
//	return database.NewConfig(dataDir, ethCfg, nodeCfg, ""), nil
//}

