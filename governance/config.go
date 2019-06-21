package governance

import (
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	contractAddress common.Address
	gethIpc         string
}

type jsonConfig struct {
	ContractAddress string `json:"contract_address"`
}

func NewConfig(address common.Address, gethIpc string) (Config) {
	return Config{
		contractAddress: address,
		gethIpc: gethIpc,
	}
}

func (c Config) ContractAddress() common.Address {
	return c.contractAddress
}