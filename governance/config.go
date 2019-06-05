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

func NewConfig(address string) (Config) {
	return Config{
		contractAddress: common.HexToAddress(address),
	}
}

func (c Config) ContractAddress() common.Address {
	return c.contractAddress
}