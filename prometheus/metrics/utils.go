package metrics

import (
	"math/big"
	"math"
)

func Web3FromWei(amountInWei *big.Int) *big.Float {
	amount := new(big.Float)
	amount.SetString(amountInWei.String())
	return new(big.Float).Quo(amount, big.NewFloat(math.Pow10(18)))
}
