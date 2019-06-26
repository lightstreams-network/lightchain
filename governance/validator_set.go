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
	"encoding/hex"
	"github.com/tendermint/tendermint/crypto"
)

const deployContractGasLimit = 2000000

type ValidatorSet struct {
	contractAddress common.Address
}

func NewValidatorSet(contractAddress common.Address) ValidatorSet {
	return ValidatorSet{
		contractAddress: contractAddress,
	}
}

func DeployContract(client *ethclient.Client, txAuth authy.Auth) (common.Address, error) {
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

func (v ValidatorSet) AddValidator(client *ethclient.Client, txAuth authy.Auth, pubKey string, address common.Address) (error) {
	contractInstance, err := bindings.NewValidatorSet(v.contractAddress, client)
	if err != nil {
		return err
	}

	ctx := context.Background()
	cfg := txclient.NewTxDefaultConfig(deployContractGasLimit)
	txOps, err := txclient.GenerateTxOpts(ctx, client, txAuth, cfg)
	if err != nil {
		return err
	}

	tx, err := contractInstance.AddValidator(txOps, convertToBytes(pubKey), address)
	if err != nil {
		return err
	}

	_, err = txclient.FetchReceipt(client, tx, cfg)
	if err != nil {
		return err
	}

	return nil
}

func (v ValidatorSet) RemoveValidator(client *ethclient.Client, txAuth authy.Auth, pubKey string, address common.Address) (error) {
	contractInstance, err := bindings.NewValidatorSet(v.contractAddress, client)
	if err != nil {
		return err
	}

	ctx := context.Background()
	cfg := txclient.NewTxDefaultConfig(deployContractGasLimit)
	txOps, err := txclient.GenerateTxOpts(ctx, client, txAuth, cfg)
	if err != nil {
		return err
	}

	tx, err := contractInstance.RemoveValidator(txOps, convertToBytes(pubKey), address)
	if err != nil {
		return err
	}

	_, err = txclient.FetchReceipt(client, tx, cfg)
	if err != nil {
		return err
	}

	return nil
}

func (v ValidatorSet) FetchPubKeySet(client *ethclient.Client, address common.Address) ([]string, error) {
	contractInstance, err := bindings.NewValidatorSet(v.contractAddress, client)
	if err != nil {
		return nil, err
	}

	validatorSetSizeBN, err := contractInstance.ValidatorSetSize(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	validatorSetSize := int(validatorSetSizeBN.Int64())
	var validatorSet []string
	for i := 0; i < validatorSetSize; i++ {
		validatorPubKey, err := contractInstance.ValidatorPubKey(&bind.CallOpts{}, big.NewInt(int64(i)))
		if err != nil {
			return nil, err
		}
		validatorSet = append(validatorSet, hex.EncodeToString(validatorPubKey[:]))
	}

	return validatorSet, nil
}

func (v ValidatorSet) ValidatorAddress(client *ethclient.Client, proposer crypto.Address) (common.Address, error) {
	if v.contractAddress.String() == common.HexToAddress("0x0").String() {
		return common.Address{}, nil
	}

	contractInstance, err := bindings.NewValidatorSetCaller(v.contractAddress, client)
	if err != nil {
		return common.Address{}, err
	}

	address, err := contractInstance.ValidatorAddress(&bind.CallOpts{}, convertToBytes(proposer.String()))
	return address, err
}

func (v ValidatorSet) IsOwner(client *ethclient.Client, owner common.Address) (bool, error) {
	contractInstance, err := bindings.NewValidatorSetCaller(v.contractAddress, client)
	if err != nil {
		return false, err
	}

	isOwner, err := contractInstance.IsOwner(&bind.CallOpts{
		From: owner,
	})

	return isOwner, err
}

func convertToBytes(pubKey string) [20]byte {
	var pubKeyBytes [20]byte
	copy(pubKeyBytes[:], pubKey)
	return pubKeyBytes
}
