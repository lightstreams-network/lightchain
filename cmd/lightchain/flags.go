package main

import (
	"gopkg.in/urfave/cli.v1"
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/lightstreams-network/lightchain/utils"
)


var (
	// flags that configure the go-ethereum node
	NodeFlags = []cli.Flag{
		ethUtils.DataDirFlag,
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
		utils.TendermintRpcListenPortFlag,
		utils.TendermintP2PListenPortFlag,
		utils.ProxyListenPortFlag,
		
		utils.TargetGasLimitFlag,
		utils.TendermintAddrFlag,
		//utils.ABCIAddrFlag,
		//utils.ABCIProtocolFlag,
		utils.VerbosityFlag,
		utils.ConfigFileFlag,
		//utils.WithTendermintFlag,
	}
)
