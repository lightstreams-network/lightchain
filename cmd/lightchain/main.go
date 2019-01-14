package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"
	"github.com/spf13/cobra"
	"flag"
	
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/lightstreams-network/lightchain/log"
)

var logger = log.NewLogger()

func main() {
	lightchainCmd := LightchainCmd()

	if err := lightchainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// LightchainCmd is the main Lightstreams PoA blockchain node.
func LightchainCmd() *cobra.Command {
	var lightchainCmd = &cobra.Command{
		Use:   "lightchain",
		Short: "Lightstreams PoA blockchain node.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	lightchainCmd.AddCommand(versionCmd)
	lightchainCmd.AddCommand(docsCmd())
	lightchainCmd.AddCommand(initCmd())
	lightchainCmd.AddCommand(runCmd())
	return lightchainCmd
}

func addDefaultFlags(cmd *cobra.Command) {
	cmd.Flags().String(LogLvlFlag.Name, LogLvlFlag.Value, LogLvlFlag.Usage)

	cmd.Flags().String(DataDirFlag.Name, DataDirFlag.Value.Value, DataDirFlag.Usage)
	cmd.MarkFlagRequired(DataDirFlag.Name)	
}

func NewNodeClientCtx(dataDir string, cmd *cobra.Command) *cli.Context {
	var boolEmptyBucket bool
	var stringEmptyBucket string
	
	app := ethUtils.NewApp("0.0.0", "the lightchain command line interface")
	flagSet := flag.NewFlagSet("FakeCli", flag.ExitOnError)
	ctx := cli.NewContext(app, flagSet, nil)
	
	flagSet.StringVar(&stringEmptyBucket, ethUtils.DataDirFlag.GetName(), dataDir, ethUtils.DataDirFlag.Usage)

	rpcEnabledFlagValue, _ := cmd.Flags().GetBool(RPCEnabledFlag.GetName())
	flagSet.BoolVar(&boolEmptyBucket, RPCEnabledFlag.GetName(), rpcEnabledFlagValue, RPCEnabledFlag.Usage)
	
	rpcListenAddrValue, _ := cmd.Flags().GetString(RPCListenAddrFlag.GetName())
	flagSet.StringVar(&stringEmptyBucket, RPCListenAddrFlag.GetName(), rpcListenAddrValue, RPCListenAddrFlag.Usage)
	
	rpcPortFlagValue, _ := cmd.Flags().GetString(RPCPortFlag.GetName())
	flagSet.StringVar(&stringEmptyBucket, RPCPortFlag.GetName(), rpcPortFlagValue, RPCPortFlag.Usage)
	
	rpcApiFlag, _ := cmd.Flags().GetString(RPCApiFlag.GetName())
	flagSet.StringVar(&stringEmptyBucket, RPCApiFlag.GetName(), rpcApiFlag, RPCApiFlag.Usage)
	
	wsEnabledFlagValue, _ := cmd.Flags().GetBool(WSEnabledFlag.GetName())
	flagSet.BoolVar(&boolEmptyBucket, WSEnabledFlag.GetName(), wsEnabledFlagValue, WSEnabledFlag.Usage)
	
	wsListenAddrFlag, _ := cmd.Flags().GetString(WSListenAddrFlag.GetName())
	flagSet.StringVar(&stringEmptyBucket, WSListenAddrFlag.GetName(), wsListenAddrFlag, WSListenAddrFlag.Usage)
	
	wsPortFlag, _ := cmd.Flags().GetString(WSPortFlag.GetName())
	flagSet.StringVar(&stringEmptyBucket, WSPortFlag.GetName(), wsPortFlag, WSPortFlag.Usage)
	
	// Default values required
	flagSet.StringVar(&stringEmptyBucket, ethUtils.GCModeFlag.GetName(), ethUtils.GCModeFlag.Value, ethUtils.GCModeFlag.Usage)
	
	ctx.Set(ethUtils.GCModeFlag.GetName(), ethUtils.GCModeFlag.Value)
	
	return ctx
}
