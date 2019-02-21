# Metrics
## Prometheus Exporter

Lightchain node inner prometheus to track and export inner metrics
regarding the execution of the current node. To enable the prometheus exporter
you need to include `--prometheus` flag when you run your local node:
``
lightchain run --datadir="${HOME}/.lightchain" --prometheus
``

Then you can access to executed prometheus exported at 
[`http://localhost:26661/metrics`](http://localhost:26661/metrics)

**Sample output**
```
...
# HELP lightchain_consensus_check_txs_total_counter Checked txs total
# TYPE lightchain_consensus_check_txs_total_counter counter
lightchain_consensus_check_txs_total_counter{module="abci"} 26
# HELP lightchain_consensus_commit_block_total_counter Commited txs total
# TYPE lightchain_consensus_commit_block_total_counter counter
lightchain_consensus_commit_block_total_counter{module="abci"} 64
# HELP lightchain_consensus_deliver_txs_total_counter Delivered txs total
# TYPE lightchain_consensus_deliver_txs_total_counter counter
lightchain_consensus_deliver_txs_total_counter{module="abci"} 25
# HELP lightchain_database_broadcasted_txs_total_counter Broadcasted txs total.
# TYPE lightchain_database_broadcasted_txs_total_counter counter
lightchain_database_broadcasted_txs_total_counter 26
...
```

## Tendermint Prometheus Exporter
If you want to also launch a prometheus exporter for tendermint service
you need to enable it in `${HOME}/.lightchain/consensus/config/config.tolm`:
```
##### instrumentation configuration options #####
[instrumentation]

# When true, Prometheus metrics are served under /metrics on
# PrometheusListenAddr.
# Check out the documentation for the list of available metrics.
prometheus = true

# Address to listen for Prometheus collector(s) connections
prometheus_listen_addr = ":26660"
...
```
In this case the `tendermint` prometheus exported is going to be exposed over:
[`http://locahost:26660/metrics`](http://locahost:26660/metrics)
