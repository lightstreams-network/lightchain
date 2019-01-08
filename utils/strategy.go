package utils

import (
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
)

// MinerRewardStrategy is a mining strategy
type MinerRewardStrategy interface {
	Receiver() common.Address
}

// ValidatorsStrategy is a validator strategy
type ValidatorsStrategy interface {
	SetValidators(validators []*tmtAbciTypes.Validator)
	CollectTx(tx *ethTypes.Transaction)
	GetUpdatedValidators() []*tmtAbciTypes.Validator
}

// Strategy encompasses all available strategies
type Strategy struct {
	MinerRewardStrategy
	ValidatorsStrategy
}
