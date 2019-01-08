package config

import (
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"fmt"
	"encoding/json"
	
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/node"
)

const (
	// Client identifier to advertise over the network
	clientIdentifier = "lightchain"
	// Environment variable for home dir
	emHome = "EMHOME"
)


// The config folder name inside the lightchain's node --datadir ~/.lightchain
const DataFolderName = "lightchain"

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
	dPath := node.DefaultDataDir()

	emHome := os.Getenv(emHome)
	if emHome != "" {
		dPath = emHome
	}

	if ctx.GlobalIsSet(ethUtils.DataDirFlag.Name) {
		dPath = ctx.GlobalString(ethUtils.DataDirFlag.Name)
	}

	if dPath == "" {
		ethUtils.Fatalf("Cannot determine default data directory, please set manually (--datadir)")
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

func ConfigPath(ctx *cli.Context) string {
	return filepath.Join(ctx.GlobalString(ethUtils.DataDirFlag.Name), DataFolderName, ConfigFilename)
}
