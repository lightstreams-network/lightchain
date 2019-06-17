package governance

import (
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	contractAddress common.Address
}

type jsonConfig struct {
	ContractAddress string `json:"contract_address"`
}

func NewConfig(address common.Address) (Config) {
	return Config{
		contractAddress: address,
	}
}

func (c Config) ContractAddress() common.Address {
	return c.contractAddress
}