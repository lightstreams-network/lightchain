// nolint=lll
package utils

import (
	"gopkg.in/urfave/cli.v1"
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
		Name:  "rpc_listen_port",
		Value: TendermintRpcListenPort,
		Usage: "This is the port that lightchain will use to connect to the tendermint.",
	}
	ProxyListenPortFlag = cli.UintFlag{
		Name:  "proxy_listen_port",
		Value: ProxyListenPort,
		Usage: "This is the port that lightchain will use to connect to the tendermint.",
	}
	TendermintP2PListenPortFlag = cli.UintFlag{
		Name:  "p2p_listen_port",
		Value: TendermintP2PListenPort,
		Usage: "This is the port that lightchain will use to connect to the tendermint.",
	}
	
	
	// ----------------------------
	// ABCI Flags

	// TendermintAddrFlag is the address that lightchain will use to connect to the tendermint core node
	// #stable - 0.4.0
	TendermintAddrFlag = cli.StringFlag{
		Name:  "tendermint_addr",
		Value: "tcp://127.0.0.1:26657",
		Usage: "This is the address that lightchain will use to connect to the tendermint core node. Please provide a port.",
	}

	// ABCIAddrFlag is the address that lightchain will use to listen to incoming ABCI connections
	// #stable - 0.4.0
	ABCIAddrFlag = cli.StringFlag{
		Name:  "abci_laddr",
		Value: "tcp://0.0.0.0:26658",
		Usage: "This is the address that the ABCI server will use to listen to incoming connection from tendermint core.",
	}

	// ABCIProtocolFlag defines whether GRPC or SOCKET should be used for the ABCI connections
	// #stable - 0.4.0
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

	// WithTendermintFlag asks to start Tendermint
	// `tendermint init` and `tendermint node` when `lightchain init`
	// and `lightchain` are invoked respectively.
	WithTendermintFlag = cli.BoolFlag{
		Name: "with-tendermint",
		Usage: "If set, it will invoke `tendermint init` and `tendermint node` " +
			"when `lightchain init` and `lightchain` are invoked respectively",
	}
)
