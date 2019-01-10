package utils

import (
	"gopkg.in/urfave/cli.v1"
	"os"
	"path/filepath"
	
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
)

const (
	// Client identifier to advertise over the network
	ClientIdentifier = "lightchain"
)


const DataFolderName = "lightchain"
const KeystoreFolderName = "keystore"
const ChainDataFolderName = "chaindata"

var ProjectRootPath = filepath.Join(os.Getenv("GOPATH"), "src/github.com/lightstreams-network", "lightchain")


//func NewLightchainConfig(path string) (consensus.Config, error) {
//	buffer, err := ioutil.ReadFile(path)
//	if err != nil {
//		return consensus.Config{}, fmt.Errorf("unable to open LS's configuration file %s", path)
//	}
//
//	cfg := new(consensus.Config)
//	err = json.Unmarshal(buffer, cfg)
//	if err != nil {
//		return consensus.Config{}, err
//	}
//
//	return *cfg, nil
//}

func MakeHomeDir(ctx *cli.Context) string {
	
	dPath := DefaultDataDir()
	
	if ctx.GlobalIsSet(DataDirFlag.Name) {
		dPath = ctx.GlobalString(DataDirFlag.Name)
	}

	if dPath == "" {
		ethUtils.Fatalf("Cannot determine default data directory, please set manually (--homedir)")
	}
	
	if err := os.MkdirAll(dPath, os.ModePerm); err != nil {
		ethUtils.Fatalf("Home folder err: %v", err)
	}

	return dPath
}

func MakeDataDir(ctx *cli.Context) string {
	homeDir := MakeHomeDir(ctx)
	dataDir := filepath.Join(homeDir, DataFolderName)
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		ethUtils.Fatalf("Data folder err: %v", err)
	}
	return dataDir
}

//func ConfigPath(ctx *cli.Context) string {
//	dataDir := MakeDataDir(ctx)
//	return filepath.Join(dataDir, consensus.ConfigFilename)
//}
