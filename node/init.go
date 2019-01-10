package node

import (
	ethLog "github.com/ethereum/go-ethereum/log"
	
	"github.com/lightstreams-network/lightchain/utils"
)

func InitNode(cfg Config) error {
	ethLog.Info("Initializing DataDir", "dir", cfg.DataDir)
	if err := utils.CreatePathIfNotExists(cfg.DataDir); err != nil {
		return err
	}

	return nil
}
