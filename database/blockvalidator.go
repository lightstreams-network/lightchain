package database

import (
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

// Dummy block validator ignoring PoW, uncles and so on.
type NullBlockValidator struct{}

var _ core.Validator = NullBlockValidator{}

// ValidateBody does not validate anything.
func (NullBlockValidator) ValidateBody(*ethTypes.Block) error {
	return nil
}

// ValidateState does not validate anything.
func (NullBlockValidator) ValidateState(block, parent *ethTypes.Block, state *state.StateDB, receipts ethTypes.Receipts, usedGas uint64) error {
	return nil
}