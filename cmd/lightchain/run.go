package main

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/lightstreams-network/lightchain/node"
	"fmt"
	"os"
	"github.com/lightstreams-network/lightchain/utils"
	"github.com/lightstreams-network/lightchain/consensus"
	"github.com/lightstreams-network/lightchain/database"
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	
	tmtServer "github.com/tendermint/tendermint/abci/server"
	tmtCommon "github.com/tendermint/tmlibs/common"
	"github.com/lightstreams-network/lightchain/log"
)

func RunCmd(ctx *cli.Context) {
	abciNode := node.CreateNode(ctx)
	node.StartNode(ctx, abciNode)

	// Fetch the registered service of this type
	var ethBackend *database.Backend
	if err := abciNode.Service(&ethBackend); err != nil {
		ethUtils.Fatalf("database ethBackend service not running: %v", err)
	}

	// In-proc RPC connection so ABCI Query can be forwarded over the database rpc
	rpcClient, err := abciNode.Attach()
	if err != nil {
		ethUtils.Fatalf("Failed to attach to the inproc geth: %v", err)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create the ABCI application - in memory or persisted to disk
	ethApp, err := consensus.CreateLightchainApplication(ethBackend, rpcClient, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ethLogger := log.NewLogger()
	ethApp.SetLogger(ethLogger.With("module", "lightchain"))

	// Start the app on the ABCI server listener
	abciAddr := fmt.Sprintf("tcp://0.0.0.0:%d", ctx.GlobalInt(utils.TendermintProxyListenPortFlag.Name))
	abciProtocol := ctx.GlobalString(utils.ABCIProtocolFlag.Name)
	abciSrv, err := tmtServer.NewServer(abciAddr, abciProtocol, ethApp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	abciLogger := log.NewLogger()
	abciSrv.SetLogger(abciLogger.With("module", "node-server"))

	if err := abciSrv.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	tmtLogger := log.NewLogger()
	tmtNode, err := consensus.CreateNewNode(ctx,  tmtLogger.With("module", "tendermint"))
	if err != nil {
		fmt.Errorf("Failed to create node: %v", err)
	}
	if err := consensus.StartNode(ctx, tmtNode); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tmtCommon.TrapSignal(func() {
		if err := tmtNode.Stop(); err != nil {
			fmt.Errorf("Error stopping Tendermint service", err)
		}
		if err := abciNode.Stop(); err != nil {
			fmt.Errorf("Error stopping Geth Node", err)
		}
		if err := abciSrv.Stop(); err != nil {
			fmt.Errorf("Error stopping ABCI service", err)
		}
		os.Exit(1)
	})
}
