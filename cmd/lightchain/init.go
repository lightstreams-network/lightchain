package main

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/lightstreams-network/lightchain/node"
	"github.com/lightstreams-network/lightchain/database"
	"github.com/lightstreams-network/lightchain/consensus"
)

func InitCmd(ctx *cli.Context) error {
	if err := node.InitNode(ctx); err != nil {
		return err
	}
	
	if err := database.Init(ctx); err != nil {
		return err
	}

	if err := consensus.InitNode(ctx); err != nil {
		return err
	}
	
	return nil
}
