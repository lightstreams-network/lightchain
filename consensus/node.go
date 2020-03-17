package consensus

import (
	"fmt"
	"time"

	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/lightstreams-network/lightchain/log"
	"github.com/tendermint/tendermint/proxy"
	"github.com/lightstreams-network/lightchain/database"

	ethRpc "github.com/ethereum/go-ethereum/rpc"
	tmtLog "github.com/tendermint/tendermint/libs/log"
	tmtNode "github.com/tendermint/tendermint/node"
	tmtP2P "github.com/tendermint/tendermint/p2p"
	tmtCommon "github.com/tendermint/tendermint/libs/common"
	tmtDbm "github.com/tendermint/tendermint/libs/db"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/lightstreams-network/lightchain/consensus/metrics"
)

type Node struct {
	tendermint *tmtNode.Node
	nodeKey    *tmtP2P.NodeKey
	cfg        *Config
	logger     tmtLog.Logger
	metrics    metrics.Metrics
	dbs        map[string]tmtDbm.DB
}

func NewNode(cfg *Config, registry *prometheus.Registry) (*Node, error) {
	logger := log.NewLogger().With("engine", "consensus")

	nodeKeyFile := cfg.tendermintCfg.NodeKeyFile()
	if ! tmtCommon.FileExists(nodeKeyFile) {
		return nil, fmt.Errorf("tendermint key file '%s' does not exists", nodeKeyFile)
	}

	nodeKey, err := p2p.LoadOrGenNodeKey(nodeKeyFile)
	if err != nil {
		return nil, err
	}

	var trackedMetrics metrics.Metrics
	if cfg.metrics {
		trackedMetrics = metrics.NewMetrics(registry)
	} else {
		trackedMetrics = metrics.NewNullMetrics()
	}

	return &Node{
		nil,
		nodeKey,
		cfg,
		logger,
		trackedMetrics,
		make(map[string]tmtDbm.DB),
	}, nil
}

func (n *Node) Start(ethRPCClient *ethRpc.Client, db *database.Database) error {
	n.logger.Debug("Creating tendermint ABCI application...")
	abci, err := NewTendermintABCI(db, ethRPCClient, n.metrics)
	if err != nil {
		return err
	}

	n.logger.Debug("Creating tendermint node...")
	n.tendermint, err = tmtNode.NewNode(
		n.cfg.tendermintCfg,
		privval.LoadOrGenFilePV(n.cfg.tendermintCfg.PrivValidatorKeyFile(), n.cfg.tendermintCfg.PrivValidatorStateFile()),
		n.nodeKey,
		proxy.NewLocalClientCreator(abci),
		tmtNode.DefaultGenesisDocProviderFunc(n.cfg.tendermintCfg),
		n.defaultDBProvider,
		tmtNode.DefaultMetricsProvider(n.cfg.tendermintCfg.Instrumentation),
		n.logger,
	)
	if err != nil {
		return err
	}

	n.logger.Info("Starting tendermint node...")
	if err := n.tendermint.Start(); err != nil {
		return err
	}
	n.logger.Info("Tendermint node started")
	n.logger.Debug("Consensus node started", "nodeInfo", n.tendermint.Switch().NodeInfo())

	return nil
}

func (n *Node) IsRunning() bool {
	if n.tendermint == nil {
		return false
	}

	return n.tendermint.IsRunning();
}

func (n *Node) Stop() error {
	if n.tendermint.IsRunning() {
		n.logger.Info("Stopping tendermint node...")
		if err := n.tendermint.Stop(); err != nil {
			return err
		}
		<-n.tendermint.Quit()
		n.logger.Info("Tendermint node stopped")
	}

	// Wait a seconds to buy time for tendermint to finish node stopping events
	<-time.After(time.Second)
	
	for ctxID, db := range n.dbs {
		n.logger.Info(fmt.Sprintf("Closing database '%s'...", ctxID))
		db.Close()
		n.logger.Info(fmt.Sprintf("Database '%s' closed", ctxID))
	}

	return nil
}

func (n *Node) defaultDBProvider(ctx *tmtNode.DBContext) (tmtDbm.DB, error) {
	dbType := tmtDbm.DBBackendType(ctx.Config.DBBackend)
	n.logger.Info(fmt.Sprintf("Loading database '%s'...", ctx.ID))
	n.dbs[ctx.ID] = tmtDbm.NewDB(ctx.ID, dbType, ctx.Config.DBDir())
	return n.dbs[ctx.ID], nil
}
