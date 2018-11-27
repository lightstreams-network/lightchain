package app

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/abci/example/code"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmLog "github.com/tendermint/tendermint/libs/log"
)

var (
	stateKey        = []byte("stateKey")
	kvPairPrefixKey = []byte("kvPairKey:")
)

type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

func loadState(db dbm.DB) State {
	stateBytes := db.Get(stateKey)
	var state State
	if len(stateBytes) != 0 {
		err := json.Unmarshal(stateBytes, &state)
		if err != nil {
			panic(err)
		}
	}
	state.db = db
	return state
}

func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}

func prefixKey(key []byte) []byte {
	return append(kvPairPrefixKey, key...)
}

//---------------------------------------------------

var _ abciTypes.Application = (*KVStoreApplication)(nil)

type KVStoreApplication struct {
	abciTypes.BaseApplication
	// validator set
	ValUpdates []abciTypes.ValidatorUpdate
	state      State
	logger     tmLog.Logger
}

func NewKVStoreApplication() *KVStoreApplication {
	state := loadState(dbm.NewMemDB())
	return &KVStoreApplication{state: state}
}

// SetLogger sets the logger for the lightchain application
// #unstable
func (app *KVStoreApplication) SetLogger(log tmLog.Logger) {
	app.logger = log
}

func (app *KVStoreApplication) Info(req abciTypes.RequestInfo) (resInfo abciTypes.ResponseInfo) {
	app.logger.Info("KVStoreApplication::Info()", "data", req)
	return abciTypes.ResponseInfo{Data: fmt.Sprintf("{\"size\":%v}", app.state.Size)}
}

// tx is either "key=value" or just arbitrary bytes
func (app *KVStoreApplication) DeliverTx(tx []byte) abciTypes.ResponseDeliverTx {
	app.logger.Info("KVStoreApplication::DeliverTx()", "data", tx)
	var key, value []byte
	parts := bytes.Split(tx, []byte("="))
	if len(parts) == 2 {
		key, value = parts[0], parts[1]
	} else {
		key, value = tx, tx
	}
	app.state.db.Set(prefixKey(key), value)
	app.state.Size += 1

	tags := []cmn.KVPair{
		{Key: []byte("app.creator"), Value: []byte("jae")},
		{Key: []byte("app.key"), Value: key},
	}
	return abciTypes.ResponseDeliverTx{Code: code.CodeTypeOK, Tags: tags}
}

func (app *KVStoreApplication) CheckTx(tx []byte) abciTypes.ResponseCheckTx {
	app.logger.Info("KVStoreApplication::CheckTx()", "data", tx)
	return abciTypes.ResponseCheckTx{Code: code.CodeTypeOK}
}

func (app *KVStoreApplication) Commit() abciTypes.ResponseCommit {
	// Using a memdb - just return the big endian size of the db
	appHash := make([]byte, 8)
	app.logger.Info("KVStoreApplication::Commit()", "data", appHash)
	binary.PutVarint(appHash, app.state.Size)
	app.state.AppHash = appHash
	app.state.Height += 1
	saveState(app.state)
	return abciTypes.ResponseCommit{Data: appHash}
}

func (app *KVStoreApplication) Query(reqQuery abciTypes.RequestQuery) (resQuery abciTypes.ResponseQuery) {
	if reqQuery.Prove {
		value := app.state.db.Get(prefixKey(reqQuery.Data))
		resQuery.Index = -1 // TODO make Proof return index
		resQuery.Key = reqQuery.Data
		resQuery.Value = value
		if value != nil {
			resQuery.Log = "exists"
		} else {
			resQuery.Log = "does not exist"
		}
		return
	} else {
		value := app.state.db.Get(prefixKey(reqQuery.Data))
		resQuery.Value = value
		if value != nil {
			resQuery.Log = "exists"
		} else {
			resQuery.Log = "does not exist"
		}
		return
	}
}

// Track the block hash and header information
func (app *KVStoreApplication) BeginBlock(req abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	app.logger.Info("KVStoreApplication::BeginBlock()", "data", req)
	// reset valset changes
	app.ValUpdates = make([]abciTypes.ValidatorUpdate, 0)
	return abciTypes.ResponseBeginBlock{}
}

// Update the validator set
func (app *KVStoreApplication) EndBlock(req abciTypes.RequestEndBlock) abciTypes.ResponseEndBlock {
	app.logger.Info("KVStoreApplication::EndBlock()", "data", req)
	return abciTypes.ResponseEndBlock{ValidatorUpdates: app.ValUpdates}
}
