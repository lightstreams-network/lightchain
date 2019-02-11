package consensus

import (
	"encoding/json"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
	abciTypes "github.com/lightstreams-network/lightchain/consensus/types"
	"math/big"
	"fmt"
)

// Query for data from the application at current or past height.
func (abci *TendermintABCI) Query(query tmtAbciTypes.RequestQuery) tmtAbciTypes.ResponseQuery {
	abci.logger.Info("Querying state", "data", query)

	type jsonRequest struct {
		Method string          `json:"method"`
		ID     json.RawMessage `json:"id,omitempty"`
		Params []interface{}   `json:"params,omitempty"`
	}

	var in jsonRequest
	if err := json.Unmarshal(query.Data, &in); err != nil {
		return tmtAbciTypes.ResponseQuery{Code: uint32(abciTypes.ErrEncodingError.Code), Log: err.Error()}
	}

	var result interface{}
	if err := abci.ethRPCClient.Call(&result, in.Method, in.Params...); err != nil {
		return tmtAbciTypes.ResponseQuery{Code: uint32(abciTypes.ErrInternalError.Code), Log: err.Error()}
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return tmtAbciTypes.ResponseQuery{Code: uint32(abciTypes.ErrInternalError.Code), Log: err.Error()}
	}

	return tmtAbciTypes.ResponseQuery{Code: tmtAbciTypes.CodeTypeOK, Value: resultBytes}
}

// Info return information about the application state.
//
// Used to sync Tendermint with the application during a handshake that happens on startup.
// Tendermint expects LastBlockAppHash and LastBlockHeight to be updated during Persist,
// ensuring that Persist is never called twice for the same block height.
func (abci *TendermintABCI) Info(req tmtAbciTypes.RequestInfo) tmtAbciTypes.ResponseInfo {
	blockchain := abci.db.Ethereum().BlockChain()
	currentBlock := blockchain.CurrentBlock()
	height := currentBlock.Number()
	root := currentBlock.Root()

	abci.logger.Info("State info", "data", req, "height", height)

	// First boot-up
	if height.Cmp(big.NewInt(0)) == 0 {
		return tmtAbciTypes.ResponseInfo{
			Data:             "ABCIEthereum",
			LastBlockHeight:  height.Int64(),
			LastBlockAppHash: []byte{},
		}
	}

	return tmtAbciTypes.ResponseInfo{
		Data:             "ABCIEthereum",
		LastBlockHeight:  height.Int64(),
		LastBlockAppHash: root[:],
	}
}

// SetOption sets non-consensus critical application specific options.
//
// E.g. Key="min-fee", Value="100fermion" could set the minimum fee required
// for CheckTx (but not ExecuteTx - that would be consensus critical).
func (abci *TendermintABCI) SetOption(req tmtAbciTypes.RequestSetOption) tmtAbciTypes.ResponseSetOption {
	abci.logger.Debug(fmt.Sprintf("Setting option key '%s' value '%s'", req.Key, req.Value))

	return tmtAbciTypes.ResponseSetOption{Code: tmtAbciTypes.CodeTypeOK, Log: ""}
}