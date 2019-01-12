package main

import (
	"gopkg.in/urfave/cli.v1"
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/lightstreams-network/lightchain/utils"
)

var (
	DataDirFlag = utils.DataDirFlag
	LogLvlFlag = utils.LogLvlFlag

	RPCEnabledFlag = ethUtils.RPCEnabledFlag
	RPCListenAddrFlag = ethUtils.RPCListenAddrFlag
	RPCPortFlag = ethUtils.RPCPortFlag
	RPCApiFlag = ethUtils.RPCApiFlag
	WSEnabledFlag = ethUtils.WSEnabledFlag
	WSListenAddrFlag = ethUtils.WSListenAddrFlag
	WSPortFlag = ethUtils.WSPortFlag

	ConsensusRpcListenPortFlag = utils.TendermintRpcListenPortFlag
	ConsensusP2PListenPortFlag = utils.TendermintP2PListenPortFlag
	ConsensusProxyListenPortFlag = utils.TendermintProxyListenPortFlag
	ConsensusProxyProtocolFlag = utils.ABCIProtocolFlag
)

var (
	// flags that configure the go-ethereum node
	NodeFlags = []cli.Flag{
		//ethUtils.DataDirFlag,
		ethUtils.KeyStoreDirFlag,
		ethUtils.NoUSBFlag,
		ethUtils.NetworkIdFlag,
		// Performance tuning
		ethUtils.CacheFlag,
		ethUtils.TrieCacheGenFlag,
		ethUtils.GCModeFlag,
		// Account settings
		ethUtils.UnlockedAccountFlag,
		ethUtils.PasswordFileFlag,
		ethUtils.VMEnableDebugFlag,
		// Logging and debug settings
		ethUtils.NoCompactionFlag,
		// Gas price oracle settings
		ethUtils.GpoBlocksFlag,
		ethUtils.GpoPercentileFlag,
		// Gas Price
		//ethUtils.GasPriceFlag,
	}

	RpcFlags = []cli.Flag{
		ethUtils.RPCEnabledFlag,
		ethUtils.RPCListenAddrFlag,
		ethUtils.RPCPortFlag,
		ethUtils.RPCCORSDomainFlag,
		ethUtils.RPCVirtualHostsFlag,
		ethUtils.RPCApiFlag,
		ethUtils.IPCDisabledFlag,
		ethUtils.WSEnabledFlag,
		ethUtils.WSListenAddrFlag,
		ethUtils.WSPortFlag,
		ethUtils.WSApiFlag,
		ethUtils.WSAllowedOriginsFlag,
	}

	// flags that configure the ABCI app
	LightchainFlags = []cli.Flag{
		utils.DataDirFlag,
		utils.TendermintRpcListenPortFlag,
		utils.TendermintP2PListenPortFlag,
		utils.TendermintProxyListenPortFlag,
		
		utils.TargetGasLimitFlag,
		//utils.TendermintAddrFlag,
		//utils.ABCIAddrFlag,
		utils.ABCIProtocolFlag,
		utils.VerbosityFlag,
		utils.ConfigFileFlag,
		//utils.WithTendermintFlag,
	}
)
