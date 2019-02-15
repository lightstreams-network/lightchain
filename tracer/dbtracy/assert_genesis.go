package dbtracy

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/rawdb"
)

func (t EthDBTracer) AssertPersistedGenesisBlock(genesis core.Genesis) {
	t.logger.Infow("Tracing if ETH DB wrote a valid genesis block to disk...", "chaindata", t.chainDataDir)

	chainDb, err := ethdb.NewLDBDatabase(t.chainDataDir, 0, 0)
	if err != nil {
		t.logger.Errorw("unable to open LDB db", "err", err)
		return
	}
	defer chainDb.Close()

	head := rawdb.ReadHeadBlockHash(chainDb)
	block := rawdb.ReadBlock(chainDb, head, 0)

	stateDB, err := state.New(block.Root(), state.NewDatabase(chainDb))
	if err != nil {
		t.logger.Errorw("unable to open new stateDB", "err", err)
		return
	}

	for addr, account := range genesis.Alloc {
		persistedBalance := stateDB.GetBalance(addr)
		if persistedBalance.Cmp(account.Balance) != 0 {
			t.logger.Errorw(
				"balance defined in genesis was not properly persisted",
				"acc", addr.Hex(),
				"genesis_balance", account.Balance.String(),
				"persisted_balance", persistedBalance.String(),
			)
			continue
		}

		t.logger.Infow(
			"balance defined in genesis was properly persisted",
			"acc", addr.Hex(),
			"genesis_balance", account.Balance.String(),
			"persisted_balance", persistedBalance.String(),
		)
	}
}