package database

import (
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

// We need a block processor that just ignores PoW and uncles and so on.
//
// NullBlockProcessor does not validate anything.
//
// #unstable
type NullBlockProcessor struct{}

var _ core.Validator = NullBlockProcessor{}

// ValidateBody does not validate anything.
//
// #unstable
func (NullBlockProcessor) ValidateBody(*ethTypes.Block) error { return nil }

// ValidateState does not validate anything.
//
// #unstable
func (NullBlockProcessor) ValidateState(
	block,
	parent *ethTypes.Block,
	state *state.StateDB,
	receipts ethTypes.Receipts,
	usedGas uint64,
) error {
	return nil
}