package main

import (
	"path"
	"gopkg.in/urfave/cli.v1"
	"github.com/mitchellh/go-homedir"
	
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	ethLog "github.com/ethereum/go-ethereum/log"
)

const (
	TendermintP2PListenPort   = uint(26656)
	TendermintRpcListenPort   = uint(26657)
	TendermintProxyListenPort = uint(26658)
	TendermintProxyProtocol   = "rpc"
)

var defaultHomeDir, _ = homedir.Dir()

var (
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
