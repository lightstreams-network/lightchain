# Lightchain Tracer

**Tracer is responsible for in-depth assertions of main blockchain internals such as Consensus State, DB, Mempool and others.**

Tracer, if enabled on runtime, can also be used to collect and analyse state of an existing blockchain node in order to debug and report bugs from users running Lightchain node in real-time.

## Enable tracing

Tracing can be enabled via the `--trace` flag and the log file can be configured using the `--tracelog` flag.
```
lightchain init --datadir=/Users/enchanterio/.lightchain --standalone --trace
```

Output example:
```
tail -f /var/folders/f9/11wh9x9j31d75nxgwv2qsj3r0000gn/T/tracer.log

{"level":"info","ts":1550240215.4713638,"caller":"dbtracy/assert_genesis.go:11","msg":"Tracing if ETH DB wrote a valid genesis block to disk...","engine":"tracer","chaindata":"/Users/enchanterio/.lightchain/database/chaindata"}
{"level":"info","ts":1550240215.4742315,"caller":"dbtracy/assert_genesis.go:41","msg":"balance defined in genesis was properly persisted","engine":"tracer","acc":"0xc916Cfe5c83dD4FC3c3B0Bf2ec2d4e401782875e","genesis_balance":"300000000000000000000000000","persisted_balance":"300000000000000000000000000"}
```