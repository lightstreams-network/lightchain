package distribution

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"github.com/davecgh/go-spew/spew"
	"github.com/lightstreams-network/lightchain/authy"
	"github.com/lightstreams-network/lightchain/database/txclient"
	"context"
)

type distribution struct {
	from common.Address
	to common.Address
	amountPht *big.Int
}

func newDistribution(from common.Address, to common.Address, amount *big.Int) distribution {
	return distribution{from, to, amount}
}

func DistributeFromCsv(csvFilePath string, contractAddr common.Address) (distributionsCount int, err error) {
	distributions, err := parseCsvDistributions(csvFilePath)
	if err != nil {
		return 0, err
	}

	for i, distribution := range distributions {
		spew.Dump(i, distribution)
	}

	return len(distributions), nil
}

func Deploy(auth authy.Auth, gethIpcPath string) (common.Address, error) {
	ctx := context.Background()
	client, err := txclient.Dial(gethIpcPath)
	if err != nil {
		return common.Address{}, err
	}
	defer client.Close()

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Address{}, err
	}

	txCfg := txclient.NewTxConfig(1000000, gasPrice.Uint64(), txclient.TxReceiptTimeout, txclient.TxReceiptInterval)
	txOpts, err := txclient.GenerateTxOpts(ctx, client, auth, txCfg)
	if err != nil {
		return common.Address{}, err
	}

	addr, _, _, err := DeployDistribution(txOpts, client)
	if err != nil {
		return common.Address{}, err
	}

	return addr, nil
}