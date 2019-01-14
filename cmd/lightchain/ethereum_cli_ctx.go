package main

import (
	"flag"
	"github.com/spf13/cobra"
	"gopkg.in/urfave/cli.v1"
	"github.com/mitchellh/go-homedir"
	"path"
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	ethLog "github.com/ethereum/go-ethereum/log"
)

const (
	TendermintP2PListenPort   = uint(26656)
	TendermintRpcListenPort   = uint(26657)
	TendermintProxyListenPort = uint(26658)
	TendermintProxyProtocol   = "rpc"
)

var (
	defaultHomeDir, _ = homedir.Dir()

	DataDirFlag = ethUtils.DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: ethUtils.DirectoryString{path.Join(defaultHomeDir, "lightchain")},
	}

	LogLvlFlag = cli.StringFlag{
		Name:  "lvl",
		Usage: "Level of logging",
		Value: ethLog.LvlInfo.String(),
	}

	RPCEnabledFlag = ethUtils.RPCEnabledFlag
	RPCListenAddrFlag = ethUtils.RPCListenAddrFlag
	RPCPortFlag = ethUtils.RPCPortFlag
	RPCApiFlag = ethUtils.RPCApiFlag
	WSEnabledFlag = ethUtils.WSEnabledFlag
	WSListenAddrFlag = ethUtils.WSListenAddrFlag
	WSPortFlag = ethUtils.WSPortFlag

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
		Value: "socket",
		Usage: "socket | grpc",
	}
)

func addEthNodeFlags(cmd *cobra.Command) {
	// RPC Flags
	cmd.Flags().Bool(RPCEnabledFlag.GetName(), false, RPCEnabledFlag.Usage)
	cmd.Flags().String(RPCListenAddrFlag.GetName(), RPCListenAddrFlag.Value, RPCListenAddrFlag.Usage)
	cmd.Flags().Int(RPCPortFlag.GetName(), RPCPortFlag.Value, RPCPortFlag.Usage)
	cmd.Flags().String(RPCApiFlag.GetName(), RPCApiFlag.Value, RPCApiFlag.Usage)

	// WS Flags
	cmd.Flags().Bool(WSEnabledFlag.GetName(), false, WSEnabledFlag.Usage)
	cmd.Flags().String(WSListenAddrFlag.GetName(), WSListenAddrFlag.Value, WSListenAddrFlag.Usage)
	cmd.Flags().Int(WSPortFlag.GetName(), WSPortFlag.Value, WSPortFlag.Usage)

	// Consensus Flags
	cmd.Flags().Uint(ConsensusRpcListenPortFlag.GetName(), ConsensusRpcListenPortFlag.Value, ConsensusRpcListenPortFlag.Usage)
	cmd.Flags().Uint(ConsensusP2PListenPortFlag.GetName(), ConsensusP2PListenPortFlag.Value, ConsensusP2PListenPortFlag.Usage)
	cmd.Flags().Uint(ConsensusProxyListenPortFlag.GetName(), ConsensusProxyListenPortFlag.Value, ConsensusProxyListenPortFlag.Usage)
	cmd.Flags().String(ConsensusProxyProtocolFlag.GetName(), ConsensusProxyProtocolFlag.Value, ConsensusProxyProtocolFlag.Usage)
}

func newNodeClientCtx(dataDir string, cmd *cobra.Command) *cli.Context {
	var boolEmptyBucket bool
	var stringEmptyBucket string

	app := ethUtils.NewApp("0.0.0", "the lightchain command line interface")
	flagSet := flag.NewFlagSet("Ethereum CLI ctx", flag.ExitOnError)
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