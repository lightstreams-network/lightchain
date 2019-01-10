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

			ctx := cli.NewContext(nil, nil, nil)

			abciNode := node.CreateNode(ctx)
			node.StartNode(ctx, abciNode)

			// Fetch the registered service of this type
			var ethBackend *database.Backend
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

			// Create the ABCI application - in memory or persisted to disk
			tendermintABCI, err := consensus.NewTendermintABCI(ethBackend, rpcClient)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			ethLogger := log.NewLogger()
			tendermintABCI.SetLogger(ethLogger.With("module", "lightchain"))

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

			tmtLogger := log.NewLogger()
			tmtNode, err := consensus.CreateNewNode(ctx,  tmtLogger.With("module", "tendermint"))
			if err != nil {
				logger.Error(fmt.Errorf("failed to create consensus node: %v", err).Error())
				os.Exit(1)
			}

			if err := consensus.StartNode(ctx, tmtNode); err != nil {
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

	return runCmd
}