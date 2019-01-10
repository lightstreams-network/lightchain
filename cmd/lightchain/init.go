package main

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
	
	"gopkg.in/urfave/cli.v1"

	"github.com/ethereum/go-ethereum/eth"
	
	"github.com/lightstreams-network/lightchain/node"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/utils"
)

func initCmd() *cobra.Command {
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initializes new lightchain node according to the configured flags.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, _ := cmd.Flags().GetString(flagDataDir)

			nodeCfg := node.NewConfig(dataDir)
			if err := node.InitNode(nodeCfg); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			consensusCfg, err := createConsensusConfig(dataDir)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			if err := consensus.InitNode(consensusCfg); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			dbCfg, err := createDatabaseConfig(dataDir)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			if err := database.Init(dbCfg); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			logger.Info(fmt.Sprintf("Lightchain node successfully initialized into '%s'!", dataDir))
			os.Exit(0)
		},
	}
	
	addDefaultFlags(initCmd)
	return initCmd
}

func createDatabaseConfig(dataDir string) (database.Config, error) {
	ethCfg := eth.DefaultConfig
	nodeCfg := database.DefaultEthNodeConfig()
	
	//ethUtils.SetNodeConfig(ctx, &nodeCfg)
	database.SetNodeDefaultConfig(&nodeCfg, dataDir)
	//
 	_, err := database.NewNode(&nodeCfg)
 	if err != nil {
		return database.Config{}, err
	}
 	//ethUtils.SetEthConfig(ctx, &stack.Node, &ethCfg)
	database.SetEthDefaultConfig(&ethCfg)
 	
	return database.NewConfig(dataDir, ethCfg, nodeCfg, ""), nil
}


func createConsensusConfig(dataDir string) (consensus.Config, error) {
	//rpcPort := ctx.GlobalUint(utils.TendermintRpcListenPortFlag.Name)
	//proxyPort := ctx.GlobalUint(utils.TendermintProxyListenPortFlag.Name)
	//p2pPort := ctx.GlobalUint(utils.TendermintP2PListenPortFlag.Name)
	rpcPort := utils.TendermintRpcListenPort
	proxyPort := utils.ProxyListenPort
	p2pPort := utils.TendermintP2PListenPort
	return consensus.NewConfig(dataDir, uint16(rpcPort), uint16(proxyPort), uint16(p2pPort)), nil
}

func createNodeConfig(ctx *cli.Context) (node.Config, error) {
	dataDir := ctx.GlobalString(utils.DataDirFlag.Name)
	return node.NewConfig(dataDir), nil
}


