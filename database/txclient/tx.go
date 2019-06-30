package txclient

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"github.com/ethereum/go-ethereum/core/types"
	"context"
	"github.com/lightstreams-network/lightchain/authy"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"time"
	"github.com/ethereum/go-ethereum"
	"fmt"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
)

func GenerateTxOpts(ctx context.Context, client *ethclient.Client, auth authy.Auth, cfg TxConfig) (*bind.TransactOpts, error) {
	txOpts := bind.NewKeyedTransactor(auth.PrivKey())

	logger.Debug(fmt.Sprintf("Obtaining pending nonce for account '%s'.", txOpts.From.Hex()))
	nonce, err := client.PendingNonceAt(ctx, txOpts.From)
	if err != nil {
		return nil, err
	}
	txOpts.Nonce = big.NewInt(int64(nonce))
	logger.Debug(fmt.Sprintf("Nonce '%d' calculated for account '%s'.", nonce, txOpts.From.Hex()))

	logger.Debug(fmt.Sprintf("Gas price '%d' configured.", cfg.GasPrice()))

	txOpts.Value = big.NewInt(0)
	txOpts.GasPrice = big.NewInt(int64(cfg.GasPrice()))
	txOpts.GasLimit = cfg.GasLimit()

	balance, err := client.BalanceAt(ctx, txOpts.From, nil)
	if err != nil {
		return nil, err
	}
	logger.Debug(fmt.Sprintf("Account '%s' current balance is '%s'.", txOpts.From.Hex(), balance.String()))

	return txOpts, nil
}

func SignTransferTx(ctx context.Context, client *ethclient.Client, auth authy.Auth, to common.Address, amount *big.Int, cfg TxConfig) (*types.Transaction, error) {
	txOpts, err := GenerateTxOpts(ctx, client, auth, cfg)
	if err != nil {
		return nil, err
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	rawTx := types.NewTransaction(txOpts.Nonce.Uint64(), to, amount, txOpts.GasLimit, txOpts.GasPrice, nil)

	return types.SignTx(rawTx, types.NewEIP155Signer(chainID), auth.PrivKey())
}

// FetchReceipt periodically checks if transaction was already mined.
//
// Returns TX receipt when transaction succeeds. Error otherwise.
func FetchReceipt(client *ethclient.Client, tx *types.Transaction, cfg TxConfig) (*types.Receipt, error) {
	var start = time.Now()
	var deadline = start.Add(cfg.TxReceiptTimeout())

	logger.Info(fmt.Sprintf("Fetching TX receipt at '%s' with deadline set to '%s' with interval '%s'", start, deadline, cfg.TxReceiptInterval().String()))

	for time.Now().Before(deadline) {
		time.Sleep(cfg.TxReceiptInterval())

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			if err == ethereum.NotFound {
				logger.Info(fmt.Sprintf("Receipt for TX '%s' not found yet...", tx.Hash().Hex()))
				continue
			}

			logger.Info(err.Error())
			continue
		}

		if receipt.Status == types.ReceiptStatusFailed {
			return nil, fmt.Errorf("TX '%s' failed with status '%d' according to receipt '%s'", tx.Hash().Hex(), receipt.Status, receipt.TxHash.Hex())
		}

		return receipt, nil
	}

	return nil, fmt.Errorf("deadline for obtaining tx receipt reached")
}

func ExtractSender(tx *ethTypes.Transaction) (common.Address, error) {
	var signer ethTypes.Signer = ethTypes.FrontierSigner{}
	if tx.Protected() {
		signer = ethTypes.NewEIP155Signer(tx.ChainId())
	}
	return ethTypes.Sender(signer, tx)
}