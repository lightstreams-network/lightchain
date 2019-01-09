package main

import (
	"path/filepath"
	"gopkg.in/urfave/cli.v1"
	"os"
	
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/lightstreams-network/lightchain/utils"
)

func ResetCmd(ctx *cli.Context) error {
	dbDir := filepath.Join(utils.MakeHomeDir(ctx), "lightchain")
	if err := os.RemoveAll(dbDir); err != nil {
		ethLog.Debug("Could not reset lightchain. Failed to remove %+v", dbDir)
		return err
	}

	ethLog.Info("Successfully removed all data", "dir", dbDir)
	return nil
}
