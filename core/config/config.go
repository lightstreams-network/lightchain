package config

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
)

// The config folder name inside the lightchain's node --datadir ~/.lightchain
const CONFIG_FOLDER string = "lightchain"

// The config filename inside the CONFIG_FOLDER
const CONFIG_NAME string = "lightstreams_config.json"

type Config struct {
	DeploymentWhitelist []string `json:"deploymentWhitelist"`
}

func New(path string) (Config, error) {
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

func ReadDefaultLsConfigBlob() ([]byte) {
	return []byte(`
{
    "deploymentWhitelist": [""]
}`)
}
