package governance

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	
	"github.com/lightstreams-network/lightchain/governance/bindings"
)

type Validators struct {
	contract common.Address
	gethIpc string
}

func New(contractAddress common.Address, gethIpc string) Validators {
	return Validators {
		contract: contractAddress,
		gethIpc: gethIpc,
	}
}

func (v Validators) ValidatorAddress(pubKey string) (common.Address, error) {
	client, err := ethclient.Dial(v.gethIpc)
	if err != nil {
		return common.Address{}, err
	}
	defer client.Close()

	contractInstance, err := bindings.NewValidators(v.contract, client)
	if err != nil {
		return common.Address{}, err
	}

	address, err := contractInstance.ValidatorAddress(&bind.CallOpts{}, pubKey)
	return address, err
}