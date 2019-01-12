// nolint=lll
package utils

import (
	"gopkg.in/urfave/cli.v1"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/cmd/utils"
)

const (
	GenesisTargetGasLimit     = uint64(100000000)
	TendermintP2PListenPort   = uint(26656)
	TendermintRpcListenPort   = uint(26657)
	TendermintProxyListenPort = uint(26658)
	TendermintProxyProtocol   = "rpc"
)

var (
	// ----------------------------
	// New Ones
	TendermintRpcListenPortFlag = cli.UintFlag{
		Name:  "tmt_rpc_port",
		Value: TendermintRpcListenPort,
		Usage: "Tendermint RPC port used to receive incoming messages from Lightchain",
	}
	TendermintProxyListenPortFlag = cli.UintFlag{
		Name:  "tmt_proxy_port",
		Value: TendermintProxyListenPort,
		Usage: "Lightchain RPC port used to receive incoming messages from Tendermint",
	}
	TendermintP2PListenPortFlag = cli.UintFlag{
		Name:  "tmt_p2p_port",
		Value: TendermintP2PListenPort,
		Usage: "Tendermint port used to achieve exchange messages across nodes",
	}
	DataDirFlag = utils.DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: utils.DirectoryString{DefaultDataDir()},
	}
	LogLvlFlag = cli.StringFlag{
		Name:  "lvl",
		Usage: "Level of logging",
		Value: ethLog.LvlInfo.String(),
	}
	GenesisPathFlag = utils.DirectoryFlag{
		Name:  "genesis",
		Usage: "Genesis path",
	}
	
	
	// ----------------------------
	// ABCI Flags

	// ABCIProtocolFlag defines whether GRPC or SOCKET should be used for the ABCI connections
	ABCIProtocolFlag = cli.StringFlag{
		Name:  "abci_protocol",
		Value: "socket",
		Usage: "socket | grpc",
	}

	// VerbosityFlag defines the verbosity of the logging
	// #unstable
	VerbosityFlag = cli.IntFlag{
		Name:  "verbosity",
		Value: 3,
		Usage: "Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=core, 5=debug, 6=detail",
	}

	// ConfigFileFlag defines the path to a TOML config for go-ethereum
	// #unstable
	ConfigFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}

	// TargetGasLimitFlag defines gas limit of the Genesis block
	// #unstable
	TargetGasLimitFlag = cli.Uint64Flag{
		Name:  "targetgaslimit",
		Usage: "Target gas limit sets the artificial target gas floor for the blocks to mine",
		Value: GenesisTargetGasLimit,
	}
)
