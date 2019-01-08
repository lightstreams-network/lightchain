package utils

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
)

var blankGenesis = new(core.Genesis)

var errBlankGenesis = errors.New("could not parse a valid/non-blank Genesis")

// parseGenesisOrDefault tries to read the content from provided
// genesisPath. If the path is empty or doesn't exist, it will
// use defaultGenesisBytes as the fallback genesis source. Otherwise,
// it will open that path and if it encounters an error that doesn't
// satisfy os.IsNotExist, it returns that error.
func ParseBlobGenesis(genesisBlob []byte) (*core.Genesis, error) {
	genesis := new(core.Genesis)
	if err := json.Unmarshal(genesisBlob, genesis); err != nil {
		return nil, err
	}

	if reflect.DeepEqual(blankGenesis, genesis) {
		return nil, errBlankGenesis
	}

	return genesis, nil
}

// Extracts sender from ethereum transaction
func ExtractSender(tx *types.Transaction) (common.Address, error) {
	var signer types.Signer = types.FrontierSigner{}
	if tx.Protected() {
		signer = types.NewEIP155Signer(tx.ChainId())
	}
 	// Make sure the transaction is signed properly
	return types.Sender(signer, tx)
}
