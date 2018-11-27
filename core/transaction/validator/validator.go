package validator

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lightstreams-network/lightchain/core/config"
	"github.com/lightstreams-network/lightchain/core/transaction/validator/deployment"
)

type Validator struct {
	deploymentValidator deployment.Validator
}

func New(cfg config.Config) (Validator) {
	return Validator{deployment.New(cfg)}
}

func (v Validator) IsValid(tx types.Transaction) (bool, error) {
	if (tx.To() == nil) {
		return v.deploymentValidator.IsValid(tx)
	}

	return true, nil
}
