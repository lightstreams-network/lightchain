package prometheus

import (
	"net/http"
	"context"
	tmtLog "github.com/tendermint/tendermint/libs/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	lcLog "github.com/lightstreams-network/lightchain/log"
)

type Node struct {
	cfg             Config
	httpSrv         *http.Server
	metricsProvider MetricsProvider
	registry        *prometheus.Registry
	logger          tmtLog.Logger
}

func NewNode(cfg Config) *Node {
	logger := lcLog.NewLogger().With("service", "prometheus")

	var metricsProvider MetricsProvider
	registry := prometheus.NewPedanticRegistry()
	
	if cfg.enabled {
		metricsProvider = NewMetricsProvider(registry, cfg.namespace, cfg.ethDialUrl)
	} else {
		metricsProvider = NewNullMetricsProvider()
	}

	return &Node{
		cfg:             cfg,
		httpSrv:         nil,
		metricsProvider: metricsProvider,
		logger:          logger,
		registry:        registry,
	}
}

func (n *Node) Start() error {
	if ! n.cfg.enabled {
		n.logger.Info("Ignored initialization of prometheus service")
		return nil
	}
	
	n.httpSrv = &http.Server{
		Addr: n.cfg.http.Addr,
		Handler: promhttp.HandlerFor(
			n.registry,
			promhttp.HandlerOpts{MaxRequestsInFlight: 3},
		),
		ReadTimeout:  n.cfg.http.ReadTimeout,
		WriteTimeout: n.cfg.http.WriteTimeout,
	}

	n.logger.Info("Prometheus endpoint opened", "addr", n.cfg.http.Addr)
	if err := n.httpSrv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (n *Node) Stop() error {
	if ! n.cfg.enabled {
		n.logger.Info("Ignored stopping of prometheus service")
		return nil
	}

	if err := n.httpSrv.Shutdown(context.Background()); err != nil {
		return err
	}

	return nil
}

func (n *Node) MetricProvider() MetricsProvider {
	return n.metricsProvider
}

//func (n *Node) MetricsHandler() http.Handler {
//	return promhttp.HandlerFor(n.registry, promhttp.HandlerOpts{
//		ErrorLog:      log.New(os.Stderr, log.Prefix(), log.Flags()),
//		ErrorHandling: promhttp.ContinueOnError,
//	})
//}
