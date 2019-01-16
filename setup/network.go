package setup

import (
	"fmt"
	
	tmtCommon "github.com/tendermint/tendermint/libs/common"
	tmtime "github.com/tendermint/tendermint/types/time"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	
	siriusConsensus "github.com/lightstreams-network/lightchain/setup/sirius/consensus"
	standaloneConsensus "github.com/lightstreams-network/lightchain/setup/standalone/consensus"
	
	siriusDatabase "github.com/lightstreams-network/lightchain/setup/sirius/database"
	standaloneDatabase "github.com/lightstreams-network/lightchain/setup/standalone/database"
)

var cdc = amino.NewCodec()

func init() {
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoRoute, nil)
}

// Name represents name of blockchain used when running a node.
type Network string

const SiriusNetwork Network = "sirius"
const StandaloneNetwork Network = "standalone"


/******************
 CONSENSUS SETUP
******************/

func ReadStandaloneConsensusConfig() ([]byte, error) {
	return []byte(standaloneConsensus.ConfigToml), nil
}

func CreateStandaloneConsensusGenesis(pv *privval.FilePV) ([]byte, error) {
	genDoc := types.GenesisDoc{
		ChainID:         fmt.Sprintf("test-chain-%v", tmtCommon.RandStr(6)),
		GenesisTime:     tmtime.Now(),
		ConsensusParams: types.DefaultConsensusParams(),
	}
	genDoc.Validators = []types.GenesisValidator{{
		Address: pv.GetPubKey().Address(),
		PubKey:  pv.GetPubKey(),
		Power:   10,
	}}
	
	genDocBytes, err := cdc.MarshalJSONIndent(genDoc, "", "  ")
	if err != nil {
		return nil, err
	}

	return genDocBytes, nil
}

func ReadSiriusConsensusGenesis() ([]byte, error) {
	return []byte(siriusConsensus.Genesis), nil
}

func ReadSiriusConsensusConfig() ([]byte, error) {
	return []byte(siriusConsensus.ConfigToml), nil
}

/******************
 DATABASE SETUP
******************/

func ReadStandaloneDatabaseGenesis() ([]byte, error) {
	return []byte(standaloneDatabase.Genesis), nil
}

func ReadStandaloneDatabaseKeystore() (map[string][]byte, error) {
	var files = make(map[string][]byte)
	for name, keystore := range standaloneDatabase.Keystore {
		files[name] = []byte(keystore)
	}

	return files, nil
}

func ReadSiriusDatabaseGenesis() ([]byte, error) {
	return []byte(siriusDatabase.Genesis), nil
}

func ReadSiriusDatabaseKeystore() (map[string][]byte, error) {
	var files = make(map[string][]byte)
	for name, keystore := range siriusDatabase.Keystore {
		files[name] = []byte(keystore)
	}

	return files, nil
}


