package collectors

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lightstreams-network/lightchain/prometheus/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ethereum/go-ethereum/common"
)

type EthGenesisBalance struct {
	ethDialUrl string
	desc       *prometheus.Desc
}

const SiriusGenesisAccount = "0x61f3b85929c20980176d4a09771284625685c40e"
const SiriusTestAccount = "0xd119b8b038d3a67d34ca1d46e1898881626a082b"

var (
	GenesisEthAccounts = []string{SiriusGenesisAccount, SiriusTestAccount}
)

func NewEthGenesisBalance(ethDialUrl string, namespace string) *EthGenesisBalance {
	return &EthGenesisBalance{
		ethDialUrl: ethDialUrl,
		desc: prometheus.NewDesc(
			fmt.Sprintf("%s_%s_eth_wallet_balance", namespace, ethereumMetricsSubsystem),
			"the balance of ethereum accounts",
			[]string{"account"},
			nil,
		),
	}
}

func (collector *EthGenesisBalance) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EthGenesisBalance) Collect(ch chan<- prometheus.Metric) {
	ethClient, err := ethclient.Dial(collector.ethDialUrl)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
	}

	for _, acc := range (GenesisEthAccounts) {
		balance, err := ethClient.PendingBalanceAt(context.Background(), common.HexToAddress(acc))
		if err != nil {
			ch <- prometheus.NewInvalidMetric(collector.desc, err)
		} else {
			phtBalance, _ := utils.Web3FromWei(balance).Float64()
			ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, phtBalance, acc)
		}
	}
}
