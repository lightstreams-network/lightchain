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

## Run blockchain usage simulation with tracing enabled

A simulate cmd executes as follows:
- initializes a new node using the `init` cmd on (only) `--standalone` ntw
- boots-up the node using the `run` cmd
- simulates a funds transfer TX
- closes the node
- re-opens the blockchain database and reconstructs the state
- asserts correct hashes, balances, nonces, gas price etc

Simulation can be run via:
```
lightchain simulate --datadir=${HOME}/.lightchain --standalone --trace
```

Output example:
```
{"level":"info","ts":1551104932.2057989,"caller":"database/tracer.go:109","msg":"Tracing if ETH DB is in a valid state after simulation...","engine":"Tracer","chaindata":"/Users/enchanterio/.lightchain/database/chaindata"}
{"level":"info","ts":1551104932.2084565,"caller":"database/tracer.go:127","msg":"correct simulated TX hash","engine":"Tracer","hash":"0x49828264ee35226d1242893343c29537a04f61bbb2671d58d7ff9e16f5e5c4bd"}
{"level":"info","ts":1551104932.2084877,"caller":"database/tracer.go:135","msg":"correct root hash","engine":"Tracer","hash":"0x39e5006078070400d77b2494c879d5a6315ab3a7ae19c1309c0b4126f1f56f9c"}
{"level":"info","ts":1551104932.2085102,"caller":"database/tracer.go:143","msg":"correct coinbase","engine":"Tracer","acc":"0x0000000000000000000000000000000000000000"}
{"level":"info","ts":1551104932.2085276,"caller":"database/tracer.go:151","msg":"correct parent hash","engine":"Tracer","hash":"0x55e06fc7b51b31efb053f128068be8b09c86569895d98591ea5790b683770c58"}
{"level":"info","ts":1551104932.2085927,"caller":"database/tracer.go:174","msg":"sender balance after TX simulation is correct","engine":"Tracer","acc":"0xc916Cfe5c83dD4FC3c3B0Bf2ec2d4e401782875e","genesis_balance":"300000000000000000000000000","expected_balance":"299999998999979000000000000","post_simulation_balance":"299999998999979000000000000"}
{"level":"info","ts":1551104932.2086222,"caller":"database/tracer.go:191","msg":"TX gas price set to default gas price as expected","engine":"Tracer","default_gas_price":"1000000000","tx_gas_price":"1000000000"}
{"level":"info","ts":1551104932.2086434,"caller":"database/tracer.go:207","msg":"correct sender nonce","engine":"Tracer","expected_nonce":1,"actual_nonce":1}
```