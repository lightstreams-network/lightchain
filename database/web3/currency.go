package web3

import (
	"math/big"
	"fmt"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/common/math"
)

func ParseWei(valueInHexOrDecimal string) (*big.Int, error) {
	amount, ok := math.ParseBig256(valueInHexOrDecimal)
	if !ok {
		return nil, fmt.Errorf("unable to convert '%v' Wei value to a big.Int", valueInHexOrDecimal)
	}

	return amount, nil
}

func WeiToPhoton(amountInWei *big.Int) *big.Float {
	amount := new(big.Float)
	amount.SetString(amountInWei.String())
	return new(big.Float).Quo(amount, big.NewFloat(params.Ether))
}

func PhtToWei(amount uint64) *big.Int {
	pht := new(big.Int)
	pht.SetUint64(amount)

	return new(big.Int).Mul(pht, big.NewInt(params.Ether))
}