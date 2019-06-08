package governance

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lightstreams-network/lightchain/authy"
	
	"github.com/lightstreams-network/lightchain/governance/bindings"
	"github.com/lightstreams-network/lightchain/database/txclient"
	"context"
	"math/big"
)

const deployContractGasLimit = 2000000

type ValidatorSet struct {
	contract common.Address
	gethIpc string
}

func NewValidatorSet(contractAddress common.Address, gethIpc string) ValidatorSet {
	return ValidatorSet{
		contract: contractAddress,
		gethIpc: gethIpc,
	}
}

func DeployContract(txAuth authy.Auth, gethIpc string) (common.Address, error) {
	client, err := ethclient.Dial(gethIpc)
	if err != nil {
		return common.Address{}, err
	}
	defer client.Close()
	
	ctx := context.Background()
	cfg := txclient.NewTxDefaultConfig(deployContractGasLimit)

	txOps, err := txclient.GenerateTxOpts(ctx, client, txAuth, cfg)
	if err != nil {
		return common.Address{}, err
	}

	addr, _, _, err := bindings.DeployValidatorSet(txOps, client)
	if err != nil {
		return common.Address{}, err
	}

	return addr, nil
}

func (v ValidatorSet) AddValidator(txAuth authy.Auth, pubKey string, address common.Address) (error) {
	client, err := ethclient.Dial(v.gethIpc)
	if err != nil {
		return err
	}

	defer client.Close()

	contractInstance, err := bindings.NewValidatorSet(v.contract, client)
	if err != nil {
		return err
	}

	ctx := context.Background()
	cfg := txclient.NewTxDefaultConfig(deployContractGasLimit)
	txOps, err := txclient.GenerateTxOpts(ctx, client, txAuth, cfg)
	if err != nil {
		return err
	}

	tx, err := contractInstance.AddValidator(txOps, pubKey, address)
	if err != nil {
		return err
	}
	
	_, err = txclient.FetchReceipt(client, tx, cfg)
	if err != nil {
		return err
	}

	return nil
}

func (v ValidatorSet) RemoveValidator(txAuth authy.Auth, pubKey string, address common.Address) (error) {
	client, err := ethclient.Dial(v.gethIpc)
	if err != nil {
		return err
	}

	defer client.Close()

	contractInstance, err := bindings.NewValidatorSet(v.contract, client)
	if err != nil {
		return err
	}

	ctx := context.Background()
	cfg := txclient.NewTxDefaultConfig(deployContractGasLimit)
	txOps, err := txclient.GenerateTxOpts(ctx, client, txAuth, cfg)
	if err != nil {
		return err
	}

	tx, err := contractInstance.RemoveValidator(txOps, pubKey, address)
	if err != nil {
		return err
	}
	
	_, err = txclient.FetchReceipt(client, tx, cfg)
	if err != nil {
		return err
	}

	return nil
}

func (v ValidatorSet) FetchPubKeySet(address common.Address) ([]string, error) {
	client, err := ethclient.Dial(v.gethIpc)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	contractInstance, err := bindings.NewValidatorSet(v.contract, client)
	if err != nil {
		return nil, err
	}

	validatorSetSizeBN, err := contractInstance.ValidatorSetSize(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	validatorSetSize :=int(validatorSetSizeBN.Int64())
	var validatorSet []string
	for i := 0; i < validatorSetSize; i++ {
		validatorPubKey, err := contractInstance.ValidatorPubKey(&bind.CallOpts{}, big.NewInt(int64(i)))
		if err != nil {
			return nil, err
		}
		validatorSet = append(validatorSet, validatorPubKey)
	}

	return validatorSet, nil
}

func (v ValidatorSet) ValidatorAddress(pubKey string) (common.Address, error) {
	client, err := ethclient.Dial(v.gethIpc)
	if err != nil {
		return common.Address{}, err
	}
	defer client.Close()

	contractInstance, err := bindings.NewValidatorSetCaller(v.contract, client)
	if err != nil {
		return common.Address{}, err
	}

	address, err := contractInstance.ValidatorAddress(&bind.CallOpts{}, pubKey)
	return address, err
}