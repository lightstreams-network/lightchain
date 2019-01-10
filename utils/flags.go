// nolint=lll
package utils

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/ethereum/go-ethereum/cmd/utils"
)

// @TODO Move to consts remaining default values
var (
	GenesisTargetGasLimit   = uint64(100000000)
	TendermintP2PListenPort = uint(26656)
	TendermintRpcListenPort = uint(26657)
	ProxyListenPort 		= uint(26658)
)


var (
	// ----------------------------
	// New Ones
	TendermintRpcListenPortFlag = cli.UintFlag{
		Name:  "tmt_rpc_port",
		Value: TendermintRpcListenPort,
		Usage: "This is the port that lightchain will use to connect to the tendermint.",
	}
	TendermintProxyListenPortFlag = cli.UintFlag{
		Name:  "tmt_proxy_port",
		Value: ProxyListenPort,
		Usage: "This is the port that tendermint will use to connect to the lightchain.",
	}
	TendermintP2PListenPortFlag = cli.UintFlag{
		Name:  "tmt_p2p_port",
		Value: TendermintP2PListenPort,
		Usage: "This is the port that tendermint nodes will use to connect achieve consensus.",
	}
	DataDirFlag = utils.DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: utils.DirectoryString{DefaultDataDir()},
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
