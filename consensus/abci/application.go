package abci

import (
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/ethereum/go-ethereum/rpc"
	tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
)

// Application is the main hook of application layer (blockchain) for connecting to consensus (Tendermint) using ABCI.
//
// Tendermint ABCI requires three connections to the application to handle the different message types.
// Connections:
//    Consensus Connection - InitChain, BeginBlock, DeliverTx, EndBlock, Commit
//    Mempool Connection - CheckTx
//    Info Connection - Info, SetOption, Query
//
// Flow:
//		1. BeginBlock
//		2. CheckTx
//	    3. DeliverTx
//		4. EndBlock
//		5. Commit
//		6. CheckTx (clean mempool from TXs not included in committed block)
//
// Tendermint runs CheckTx and DeliverTx concurrently with each other,
// though on distinct ABCI connections - the mempool connection and the consensus connection.
//
// Full ABCI specs: https://tendermint.com/docs/spec/abci/abci.html
type Application struct {
	*consensusConnection
	*mempoolConnection
	*infoConnection
}

var _ tmtAbciTypes.Application = &Application{}

// NewApplication creates a new instance of Tendermint ABCI with clean block's state.
func NewApplication(db *database.Database, ethRPCClient *rpc.Client) (*Application, error) {
	logger := log.NewLogger().With("engine", "consensus", "module", "ABCI")

	ethState, err := db.Ethereum().BlockChain().State()
	if err != nil {
		return nil, err
	}

	mempoolConnection := newMempoolConnection(ethState.Copy(), logger)

	consConnection, err := newConsensusConnection(db, mempoolConnection.replaceEthState, logger)
	if err != nil {
		return nil, err
	}

	infoConnection := newInfoConnection(
		db.Ethereum().BlockChain().CurrentBlock,
		ethRPCClient,
		logger,
	)

	app := &Application{consConnection, mempoolConnection, infoConnection}

	err = app.ResetBlockState()
	if err != nil {
		return nil, err
	}

	return app, nil
}