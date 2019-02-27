# Lightchain

This is the official Lightstreams implementation of a proof-of-authority (PoA) blockchain. Lightchain is an ethereum-compatible blockchain which uses byzantine consensus to replace the original proof-of-work (PoW) from Ethereum. This is achieved by integrating [`tendermint`](https://github.com/tendermint/tendermint) for the consensus layer.

## About Lightstreams

We are currently working hard to release the **Lightstreams main network** which aims to provide a fast, privacy-enabled, content-sharing blockchain. Stay tuned about our project's progress by reading [the lightstreams blog](https://medium.com/lightstreams) or by checking out the [Lightstreams' website](https://www.lightstreams.network).

## Documentation

You can find more detailed documentation in the [lightchain CLI reference documentation](https://docs.lightstreams.network/cli-docs/lightchain).

## Pre-requirements

- Go >= 1.10
- dep (Go dependence manager) [Installation guide](https://golang.github.io/dep/docs/installation.html)

## Installation

To install `lightchain` in your system just run following commands:
```
mkdir -p ${GOPATH}/src/github.com/lightstreams-network
cd ${GOPATH}/src/github.com/lightstreams-network
git clone https://github.com/lightstreams-network/lightchain.git ./lightchain
cd ./lightchain
make get_vendor_deps
make install
```

To validate if you installation was done correctly, run the following command to obtain current installed version of Lightchain
```
lightchain version
```

It should output something like this
```
Version: 0.9.1-alpha Sirius-Net
```

## Create a Lightstreams node

Lightstreams provides a testnet called `sirius`. By default, all new created nodes get connected to this network and are automatically synchronized. **Please note** that we are actively working on improving the performance and stability of the network, therefore some issues might still occur which force us to restore blockchain. 

### Node initialization

To initialise a new blockchain you need to run `lightchain init` and  choose a local path where blockchain files are going to be stored.
```
lightchain init --datadir="${HOME}/.lightchain"
```

### Node launch

To run a lightchain node you only need to run the following command:
```
lightchain run --datadir="${HOME}/.lightchain"
```

After you run the command above, the network synchronization will take several minutes. So grab a coffee and [request some test tokens](https://discuss.lightstreams.network/t/request-test-tokens/64) while you wait :)

### Node launch with RPC open

*Note: running Lightchain node with RPC open can be dangerous due to same reasons like when using `Geth`!*

To run a lightchain node with RPC open, you only need to append the RPC flags as in Geth, final command:
```
lightchain run --datadir="${HOME}/.lightchain" --rpc --rpcaddr=0.0.0.0 --rpcport=8545 --rpcapi eth,net,web3,personal,admin
```

### Request free tokens
Please sign up to the [Lightstreams Community Forum](https://discuss.lightstreams.network) and [follow the instructions in this thread](https://discuss.lightstreams.network/t/request-test-tokens/64) to obtain free tokens to use in our test network Sirius.

### Block explorer
To see the current state of the `sirius` network and check the status of your transactions, you can go to the [lightstreams block explorer](https://explorer.sirius.lightstreams.io/home)

#### Available flags

When you run `lightchain run` or `lightchain run --help`, you will see a list of available flags:

```
Launches lightchain node and all of its online services including blockchain (Geth) and the consensus (Tendermint).

Usage:
  lightchain run [flags]

Flags:
      --abci_protocol string   socket | grpc (default "socket")
      --datadir string         Data directory for the databases and keystore (default "/home/a/lightchain")
  -h, --help                   help for run
      --prometheus             Enable prometheus metrics exporter
      --lvl string             Level of logging (default "info")
      --rpc                    Enable the HTTP-RPC server
      --rpcaddr string         HTTP-RPC server listening interface (default "localhost")
      --rpcapi string          API's offered over the HTTP-RPC interface
      --rpcport int            HTTP-RPC server listening port (default 8545)
      --tmt_p2p_port uint      Tendermint port used to achieve exchange messages across nodes (default 26656)
      --tmt_proxy_port uint    Lightchain RPC port used to receive incoming messages from Tendermint (default 26658)
      --tmt_rpc_port uint      Tendermint RPC port used to receive incoming messages from Lightchain (default 26657)
      --ws                     Enable the WS-RPC server
      --wsaddr string          WS-RPC server listening interface (default "localhost")
      --wsport int             WS-RPC server listening port (default 8546)

```

#### Run in standalone mode
Standalone mode allows you to create an isolated node for testing proposes. To do it, you can run the following command:

```
lightchain init --datadir="${HOME}/.lightchain" --standalone
```

At the genesis block, the account `0xc916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e`has been initialized with _300M Photons_. The passphrase is `WelcomeToSirius`

## Metrics
[Read how to run metrics explorer](METRICS.md)

## Troubleshooting

### Corrupted Database state

If you node is displaying the following error message:
```
Nonce not strictly increasing. Expected YYYY got XXXX engine=consensus module=ABCI
```

That was likely caused due to an synchronization issue after a transaction
 was executed from your local node. Due to that connection problem the 
 local state becomes invalid and the nonce of account used to performes
 that transaction stays behind the real one.
 
 We have an open issue to resolve this problem 
 
 In meanwhile our team resolve the issue [#70](https://github.com/lightstreams-network/lightchain/issues/70), 
 there is an alternative solution: 
 * Shut down your node
 * Remove the memory state saved on `$rm ${DATADIR}/database/transactions.rlp`
 * Start the node again `lightchain run...`

## Applications

### Leth
Lightstreams created a command line application called `leth` to run and manage a lightstreams node, as well as interact with the Lightstreams network
 - [leth documentation](https://docs.lightstreams.network/getting-started/).

Leth wraps `geth` & `ipfs` into a simple, easy-to-use interface and which connects to the Lightstreams network. For convenience, we provide an [HTTP gateway API](https://docs.lightstreams.network/api-docs/) or you can also use the [command line client](https://docs.lightstreams.network/cli-docs/leth/)

## Docker
In case you prefer to use Docker, follow the instructions below.

First, create a new docker image, which will be tagged as `lightchain:latest`
```
make docker
```

Once that is completed, you just need to run the following command
which will create your container with a running instance of lightchain
```
docker run -p 8545:8545 -p 26657:26657 -p 26656:26656 -it lightchain:latest
```

As you can see several ports has been mapped to your local environment:
- `8545` which exposes the rpc api of Ethereum
```
geth attach http://localhost:8545
```
- `26657` websocket api of tendermint
```
curl -X http://localhost:26657/status
```
- `26656` required by the consensus engine (Tendermint) for p2p communications



## Project data structure

```
├── consensus
│   ├── config
│   │   ├── config.toml
│   │   ├── genesis.json
│   │   ├── node_key.json
│   │   └── priv_validator.json
│   └── data
└── database
    ├── chaindb
    │   ├── 000001.log
    │   ├── CURRENT
    │   ├── LOCK
    │   ├── LOG
    │   └── MANIFEST-000000
    ├── genesis.json
    └── keystore
        ├── UTC--2018-08-26T21-40-07.289727986Z--4eaaad8ea38d5ef75ebdeb3d1be59d56f86c4ca9
        └── UTC--2018-08-26T21-40-31.689362077Z--4f5adedca6d869e9f5f7dcf4b7a9dfa8231a095f

```

Lightchain `datadir` is split into two main folder:
1. **consensus**: contains all the information regarding consensus
2. **database**: contains all files related to the Ethereum-compatible blockchain

## Wiki
To know more about how Lightchain works and how Tendermint is integrated to perform the PoA, visit our repository [wiki](https://github.com/lightstreams-network/lightchain/wiki)

## Tests
### Running Web3 tests verifying blockchain functionalities
[Read how to run Web3 tests](truffle/README.md)

### Running Lightchain Tracer asserting blockchain internals, state, DB
[Read how to run the Tracer](TRACER.md)

## Credit & Licenses

- go-ethereum
  - [source](https://github.com/ethereum/go-ethereum),[license](https://github.com/ethereum/go-ethereum/#license)
- tendermint
  - [source](https://github.com/tendermint/tendermint),[license](https://github.com/tendermint/tendermint/blob/master/LICENSE)

## Bugs, Issues, Questions
If you find any bugs or simply have a question, please [write an issue](https://github.com/lightstreams-network/lightchain/issues) and we'll try and help as best we can.   

## Contributors

- [Lukáš Lukáč](https://github.com/EnchanterIO), [Twitter](https://twitter.com/BlocksByLukas)
- [Gabriel Garrido](https://github.com/ggarri)
- [Andrew Zappella](https://github.com/azappella)
- [Michael Smolenski](https://github.com/mikesmo)
