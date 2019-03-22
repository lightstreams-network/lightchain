package distribution

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"github.com/lightstreams-network/lightchain/authy"
	"github.com/lightstreams-network/lightchain/database/txclient"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lightstreams-network/lightchain/log"
)

var logger = log.NewLogger().With("module", "distribution")

type newAuth = func(from common.Address) (authy.Auth, error)

type distribution struct {
	from   common.Address
	to     common.Address
	amount *big.Int
}

func newDistribution(from common.Address, to common.Address, amount *big.Int) distribution {
	return distribution{from, to, amount}
}

func DistributeFromCsv(csvFilePath string, gethIpcPath string, contractAddr common.Address, auth newAuth) (distributionsCount int, err error) {
	address2distributionsCount, distributions, err := parseCsvDistributions(csvFilePath)
	if err != nil {
		return 0, err
	}

	if len(distributions) == 0 {
		return 0, fmt.Errorf("no distributions found in CSV")
	}

	ctx := context.Background()
	client, err := txclient.Dial(gethIpcPath)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return 0, err
	}

	txCfg := txclient.NewTxConfig(100000, gasPrice.Uint64(), txclient.TxReceiptTimeout, txclient.TxReceiptInterval)

	contract, err := NewDistribution(contractAddr, client)
	if err != nil {
		return 0, err
	}

	depositTxs := make([]*types.Transaction, len(distributions))
	failedDepositTXs := make([]*types.Transaction, 1)

	from2auth := make(map[common.Address]authy.Auth)
	for from, _ := range address2distributionsCount {
		fromAuth, err := auth(from)
		if err != nil {
			return 0, err
		}

		from2auth[from] = fromAuth
	}

	for i, distribution := range distributions {
		txOpts, err := txclient.GenerateTxOpts(ctx, client, from2auth[distribution.from], txCfg)
		if err != nil {
			return 0, err
		}
		txOpts.Value = distribution.amount

		logger.Info("Depositing funds...", "from", txOpts.From.String(), "to", distribution.to, "value", txOpts.Value.String())
		tx, err := contract.Deposit(txOpts, distribution.to)
		if err != nil {
			logger.Error(err.Error())
			failedDepositTXs = append(failedDepositTXs, tx)
			continue
		}

		depositTxs[i] = tx
	}

	successfulDepositReceipts := make([]*types.Receipt, 0)
	for _, tx := range depositTxs {
		receipt, err := txclient.FetchReceipt(client, tx, txCfg)
		if err != nil {
			logger.Error(err.Error())
			failedDepositTXs = append(failedDepositTXs, tx)
			continue
		}

		if receipt.Status == 0 {
			logger.Error("deposit TX failed according to the receipt")
			failedDepositTXs = append(failedDepositTXs, tx)
			continue
		}

		logger.Info("Tokens deposited!", "to", tx.To().String(), "value", tx.Value().String(), "receipt", receipt.TxHash.String())
		successfulDepositReceipts = append(successfulDepositReceipts, receipt)
	}

	return len(successfulDepositReceipts), nil
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