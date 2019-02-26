package web3

import (
	"math/big"
	"fmt"
	"github.com/ethereum/go-ethereum/params"
)

func PhotonToWei(value string) (*big.Int, error) {
	valueBn, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return nil, fmt.Errorf("unable to convert %s into Wei", value)
	}

	return new(big.Int).Mul(valueBn, new(big.Int).Set(big.NewInt(params.Ether))), nil
}