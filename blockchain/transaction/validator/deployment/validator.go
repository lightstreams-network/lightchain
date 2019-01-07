package deployment

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lightstreams-network/lightchain/config"
	"github.com/lightstreams-network/lightchain/utils"
	"github.com/lightstreams-network/lightchain/blockchain/transaction/validator/deployment/whitelist"
)

type Validator struct {
	accountWhitelist whitelist.AccountWhitelist
}

func New(cfg config.Config) (Validator) {
	return Validator{whitelist.NewAccountWhitelist(cfg)}
}

func (v Validator) IsValid(tx types.Transaction) (bool, error) {
	sender, err := utils.ExtractSender(&tx)
	if err != nil {
		return false, err
	}

	return v.accountWhitelist.HasAccount(sender), nil
}
