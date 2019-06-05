package governance

import (
	"os"
	"path"
	"fmt"
	"github.com/lightstreams-network/lightchain/fs"
	"io/ioutil"
	"encoding/json"
	"github.com/lightstreams-network/lightchain/network"
	"github.com/ethereum/go-ethereum/common"
)

const governanceFolder = "governance"
const configFile = "config.json"

func Init(ntw network.Network, dataDir string) (Config, error) {
	governancePath := path.Join(dataDir, governanceFolder)
	if isEmpty, err := fs.IsDirEmptyOrNotExists(governancePath); !isEmpty || err != nil {
		if err != nil {
			return Config{}, err
		}
		return Config{}, fmt.Errorf("governance folder '%s' is not empty.", governancePath)
	}
	
	if err := os.MkdirAll(governancePath, os.ModePerm); err != nil {
		return Config{}, err
	}

	cfg, err := newDefaultConfig(ntw)
	if err != nil {
		return Config{}, err
	}
	
	cfgFilePath := path.Join(governancePath, configFile)
	err = writeConfigInDisk(cfgFilePath, cfg)
	if err != nil {
		return Config{}, err
	}

	return Config{}, nil
}

func Load(dataDir string) (Config, error) {
	filePath := path.Join(dataDir, governanceFolder, configFile)

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

func writeConfigInDisk(filePath string, cfg Config) error {
	jsonCfg := jsonConfig{
		ContractAddress: cfg.contractAddress.String(),
	}

	jsonFsCfg, err := json.MarshalIndent(jsonCfg, "", "	")
	if err != nil {
		return fmt.Errorf("unable to marshal governance config file to write it to the node dir '%s'. %s", filePath, err.Error())
	}

	return ioutil.WriteFile(filePath, jsonFsCfg, 0644)
}

func newDefaultConfig(ntw network.Network) (Config, error) {
	fsCfgJson, err := ntw.GovernanceConfig()
	if err != nil {
		return Config{}, fmt.Errorf("unable to read %v Governance config")
	}

	jsonCfgUnmarshal := jsonConfig{}
	err = json.Unmarshal(fsCfgJson, &jsonCfgUnmarshal)
	if err != nil {
		return Config{}, fmt.Errorf("unable to unmarshal %v Leth node's config from JSON")
	}

	return Config{
		contractAddress: common.HexToAddress(jsonCfgUnmarshal.ContractAddress),
	}, nil
}