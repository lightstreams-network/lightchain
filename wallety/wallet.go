package wallety

import (
	"math/big"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lightstreams-network/lightchain/authy"
	"github.com/lightstreams-network/lightchain/database/txclient"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
)

var logger = log.NewLogger().With("module", "wallety")

func BalanceOf(client *ethclient.Client, account common.Address) (*big.Int, error) {
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func Transfer(client *ethclient.Client, auth authy.Auth, to common.Address, valueWei string, cfg txclient.TxConfig) (*types.Transaction, error) {
	ctx := context.Background()

	amount, ok := new(big.Int).SetString(valueWei, 10)
	if !ok {
		return nil, fmt.Errorf("unable to convert '%s' Wei value to a big.Int", valueWei)
	}

	tx, err := txclient.SignTransferTx(ctx, client, auth, to, amount, cfg)
	if err != nil {
		return nil, err
	}

	err = client.SendTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	_, err = txclient.FetchReceipt(client, tx, cfg)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Account '%s' transferred '%s' Wei to account '%s'.", auth.Address().Hex(), amount, to.String()))

	return tx, nil
}