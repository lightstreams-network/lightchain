package config

import (
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"fmt"
	"encoding/json"
	
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/lightstreams-network/lightchain/utils"
)

const (
	// Client identifier to advertise over the network
	clientIdentifier = "lightchain"
	// Environment variable for home dir
	emHome = "EMHOME"
)


const DataFolderName = "lightchain"
const KeystoreFolderName = "keystore"
const ChainDataFolderName = "chaindata"

// The config filename inside the DataFolderName
const ConfigFilename = "config.json"

type Config struct {
	DeploymentWhitelist []string `json:"deploymentWhitelist"`
}

var projectRootPath = filepath.Join(os.Getenv("GOPATH"), "src/github.com/lightstreams-network", "lightchain")


func NewLightchainConfig(path string) (Config, error) {
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("unable to open LS's configuration file %s", path)
	}

	cfg := new(Config)
	err = json.Unmarshal(buffer, cfg)
	if err != nil {
		return Config{}, err
	}

	return *cfg, nil
}

func ReadDefaultConfig() ([]byte) {
	return []byte(`
{
    "deploymentWhitelist": [""]
}`)
}


func MakeHomeDir(ctx *cli.Context) string {
	
	dPath := utils.DefaultHomeDir()
	
	if ctx.GlobalIsSet(utils.HomeDirFlag.Name) {
		dPath = ctx.GlobalString(utils.HomeDirFlag.Name)
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

func MakeKeystoreDir(ctx *cli.Context) string {
	homeDir := MakeHomeDir(ctx)
	keystoreDir := filepath.Join(homeDir, KeystoreFolderName)
	if err := os.MkdirAll(keystoreDir, os.ModePerm); err != nil {
		ethUtils.Fatalf("Keystore folder err: %v", err)
	}
	return keystoreDir
}

func MakeChainDataDir(ctx *cli.Context) string {
	dataDir := MakeDataDir(ctx)
	chainDataDir := filepath.Join(dataDir, ChainDataFolderName)
	if err := os.MkdirAll(chainDataDir, os.ModePerm); err != nil {
		ethUtils.Fatalf("ChainData folder err: %v", err)
	}
	return chainDataDir
}

func ConfigPath(ctx *cli.Context) string {
	dataDir := MakeDataDir(ctx)
	return filepath.Join(dataDir, ConfigFilename)
}
