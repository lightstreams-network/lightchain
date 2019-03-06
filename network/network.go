package network

import (
	"fmt"
	
	tmtCommon "github.com/tendermint/tendermint/libs/common"
	tmtime "github.com/tendermint/tendermint/types/time"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	
	mainnetConsensus "github.com/lightstreams-network/lightchain/network/mainnet/consensus"
	siriusConsensus "github.com/lightstreams-network/lightchain/network/sirius/consensus"
	standaloneConsensus "github.com/lightstreams-network/lightchain/network/standalone/consensus"
	
	mainnetDatabase "github.com/lightstreams-network/lightchain/network/mainnet/database"
	siriusDatabase "github.com/lightstreams-network/lightchain/network/sirius/database"
	standaloneDatabase "github.com/lightstreams-network/lightchain/network/standalone/database"
	"github.com/tendermint/tendermint/version"
)

var cdc = amino.NewCodec()

func init() {
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoName, nil)
}

// Name represents name of blockchain used when running a node.
type Network string

const MainNetNetwork Network = "mainnet"
const SiriusNetwork Network = "sirius"
const StandaloneNetwork Network = "standalone"


func (n Network) ConsensusConfig() ([]byte, error) {
	switch n {
	case MainNetNetwork:
		return []byte(mainnetConsensus.ConfigToml), nil
	case SiriusNetwork:
		return []byte(siriusConsensus.Genesis), nil
	case StandaloneNetwork:
		return []byte(standaloneConsensus.ConfigToml), nil
	default:
		return []byte{}, fmt.Errorf("invalid network selected %s", n)
	}
}

func (n Network) ConsensusGenesis(pv *privval.FilePV) ([]byte, error) {
	switch n {
	case MainNetNetwork:
		return []byte(mainnetConsensus.Genesis), nil
	case SiriusNetwork:
		return []byte(siriusConsensus.Genesis), nil
	case StandaloneNetwork:
		return createConsensusGenesis(pv)
	default:
		return []byte{}, fmt.Errorf("invalid network selected %s", n)
	}
}


func (n Network) ConsensusProtocolBlockVersion() (version.Protocol, error) {
	switch n {
	case MainNetNetwork:
		return 10, nil
	case SiriusNetwork:
		return 9, nil
	case StandaloneNetwork:
		return version.BlockProtocol, nil
	default:
		return version.BlockProtocol, fmt.Errorf("invalid network selected %s", n)
	}
}

func (n Network) DatabaseGenesis() ([]byte, error) {
	switch n {
	case MainNetNetwork:
		return []byte(mainnetDatabase.Genesis), nil
	case SiriusNetwork:
		return []byte(siriusDatabase.Genesis), nil
	case StandaloneNetwork:
		return []byte(standaloneDatabase.Genesis), nil
	default:
		return []byte{}, fmt.Errorf("invalid network selected %s", n)
	}
}

func (n Network) DatabaseKeystore() (map[string][]byte, error) {
	switch n {
	case MainNetNetwork:
		var files = make(map[string][]byte)
		for name, keystore := range mainnetDatabase.Keystore {
			files[name] = []byte(keystore)
		}
	
		return files, nil
	case SiriusNetwork:
		var files = make(map[string][]byte)
		for name, keystore := range siriusDatabase.Keystore {
			files[name] = []byte(keystore)
		}
	
		return files, nil
	case StandaloneNetwork:
		var files = make(map[string][]byte)
		for name, keystore := range standaloneDatabase.Keystore {
			files[name] = []byte(keystore)
		}
	
		return files, nil
	default:
		return map[string][]byte{}, fmt.Errorf("invalid network selected %s", n)
	}
}


func createConsensusGenesis(pv *privval.FilePV) ([]byte, error) {
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
