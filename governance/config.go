package governance

import (
	"fmt"
    "io/ioutil"
    "encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"path"
	"github.com/lightstreams-network/lightchain/fs"
	"os"
)


const governanceFolder = "governance"

type Config struct {
	contractAddress common.Address
}

type jsonConfig struct {
	ContractAddress string `json:"contract_address"`
}

func NewConfig(address string) (Config) {
	return Config{
		contractAddress: common.HexToAddress(address),
	}
}

func (c Config) ContractAddress() common.Address {
	return c.contractAddress
}

func NewConfigFromDisk(dataDir string) (Config, error) {
	filePath := path.Join(dataDir, governanceFolder, "config.json")

	if !common.FileExist(filePath) {
		return Config{}, fmt.Errorf("unable to load %v Governance config", filePath)
	}

	fsCfgJson, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("unable to read %v Governance config", filePath)
	}

	jsonCfgUnmarshal := jsonConfig{}
	err = json.Unmarshal(fsCfgJson, &jsonCfgUnmarshal)
	if err != nil {
		return Config{}, fmt.Errorf("unable to unmarshal %v Leth node's config from JSON", filePath)
	}

	return Config{
		contractAddress: common.HexToAddress(jsonCfgUnmarshal.ContractAddress),
	}, nil
}

func WriteConfigInDisk(dataDir string, cfg Config) error {
	governancePath := path.Join(dataDir, governanceFolder)
	exists, err := fs.DirExists(governancePath)
	if err != nil{
		return err
	}

	if !exists {
		if err := os.MkdirAll(governancePath, os.ModePerm); err != nil {
			return err
		}
	}

	filePath := path.Join(governancePath, "config.json")
	jsonCfg := jsonConfig{
		ContractAddress: cfg.contractAddress.String(),
	}
	jsonFsCfg, err := json.MarshalIndent(jsonCfg, "", "	")
	if err != nil {
		return fmt.Errorf("unable to marshal governance config file to write it to the node dir '%s'. %s", filePath, err.Error())
	}

	return ioutil.WriteFile(filePath, jsonFsCfg, 0644)
}
