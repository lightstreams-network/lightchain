# Lightchain

This is the official Lightstreams implementation of a proof-of-authority (PoA) blockchain. Lightchain is an ethereum-compatible blockchain which uses byzantine consensus to replace the original proof-of-work (PoW) from Ethereum. This is achieved by integrating [`tendermint`](https://github.com/tendermint/tendermint) for the consensus layer.

## About Lightstreams

We are currently working hard to release the **Lightstreams main network** which aims to provide a fast, privacy-enabled, content-sharing blockchain. Stay tuned about our project's progress by reading [the lightstreams blog](https://medium.com/lightstreams) or by checking out the [Lightstreams' website](https://www.lightstreams.network).

## Documentation

You can find more detailed documentation in the [lightchain CLI reference documentation](https://docs.lightstreams.network/products-1/lightchain).

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

Note: `--> dep ensure` can take some time.

To validate if you installation was done correctly, run the following command to obtain current installed version of Lightchain
```
lightchain version
```

It should output something like this
```
Version: 1.0.0-beta Mainnet
```

## Create a Lightstreams node

Lightstreams provides its own MainNet network, and a testnet called `sirius`. 
By default, all new created nodes get connected to MainNet and synchronized, but alternatively
you can choose to connect to `sirius` or even run isolate network using `standalone`, more information is detailed
above.
 
### Node initialization

To initialise a new node you need to run `lightchain init`. The local path where blockchain 
files are going to be stored is set to `${HOME}/.lightchain`. You can change using the `--datadir` flag
```
lightchain init
```

### Node launch

To run a lightchain node you only need to run the following command:
```
lightchain run
```

After you run the command above, the network synchronization will take several minutes, so grab a coffee
while you wait :)  

### Node launch with RPC open

*Note: running Lightchain node with RPC open can be dangerous due to same reasons like when using `Geth`!*

To run a lightchain node with RPC open, you only need to append the RPC flags as in Geth, final command:
```
lightchain run --datadir="${HOME}/.lightchain" --rpc --rpcaddr=0.0.0.0 --rpcport=8545 --rpcapi eth,net,web3,personal,admin
```
## Networks

### Mainnet

To see the current state of the `MainNet` network and check the status of your transactions, 
you can go to the **[lightstreams block explorer](https://explorer.mainnet.lightstreams.io/home)**


### Sirius

To run a node using `sirius` network you will need to initialize your node using `--sirius` flag as follow:
```
lightchain init --datadir="${HOME}/.lightchain_sirius" --sirius
```

To see the current state of the `sirius` network and check the status of your transactions, 
you can go to the **[lightstreams block explorer](https://explorer.sirius.lightstreams.io/home)**

To **request free tokens** Please sign up to the [Lightstreams Community Forum](https://discuss.lightstreams.network) and [follow the instructions in this thread](https://discuss.lightstreams.network/t/request-test-tokens/64) to obtain free tokens to use in our test network Sirius.

### Standalone
Standalone mode allows you to create an isolated node for testing proposes. To do it, you can run the following command:

```
lightchain init --datadir="${HOME}/.lightchain_standalone" --standalone
```

Using this network you won't need to synchronize to anyone as you a running an independent network
where you are the only validator node.
 
To use this network at the genesis block, the account `0xc916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e` has been initialized with _300M Photons_. 
Its passphrase is `WelcomeToSirius`.

## Documentation

When you run `lightchain run` or `lightchain run --help`, you will see a list of available flags:

```
Launches lightchain node and all of its online services including blockchain (Geth) and the consensus (Tendermint).

Usage:
  lightchain run [flags]

Flags:
      --abci_protocol string   socket | grpc (default "socket")
      --datadir string         Data directory for the databases and keystore (default "/home/ggarrido/.lightchain")
  -h, --help                   help for run
      --lvl string             Level of logging (default "info")
      --prometheus             Enable prometheus metrics exporter
      --rpc                    Enable the HTTP-RPC server
      --rpcaddr string         HTTP-RPC server listening interface (default "localhost")
      --rpcapi string          API's offered over the HTTP-RPC interface
      --rpccorsdomain string   Comma separated list of domains from which to accept cross origin requests (browser enforced)
      --rpcport int            HTTP-RPC server listening port (default 8545)
      --rpcvhosts string       Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard. (default "localhost")
      --tmt_p2p_port uint      Tendermint port used to achieve exchange messages across nodes (default 26656)
      --tmt_proxy_port uint    Lightchain RPC port used to receive incoming messages from Tendermint (default 26658)
      --tmt_rpc_port uint      Tendermint RPC port used to receive incoming messages from Lightchain (default 26657)
      --trace                  Whenever to be asserting and reporting blockchain state in real-time (testing, debugging purposes)
      --tracelog string        The filepath to a log file where all tracing output will be persisted (default "/tmp/tracer.log")
      --ws                     Enable the WS-RPC server
      --wsaddr string          WS-RPC server listening interface (default "localhost")
      --wsport int             WS-RPC server listening port (default 8546)

```

If you want to know more about how to use `lightchain` command line client, see [our online documentation](https://docs.lightstreams.network/products-1/lightchain/lightchain-commands)

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
 - [leth documentation](https://docs.lightstreams.network/products-1/smart-vault/getting-started).

Leth wraps `geth` & `ipfs` into a simple, easy-to-use interface and which connects to the Lightstreams network. For convenience, we provide an 
[HTTP gateway API](https://docs-api.lightstreams.network/) or you can also use the [command line client](https://docs.lightstreams.network/products-1/smart-vault/sdk/leth)

## Docker
In case you prefer to use Docker, follow the instructions below.

First, create a new docker image, which will be tagged as `lightchain:latest`
```
make docker
```
or in case you prefer to fetch the latest version published on AWS, tagged as `lightchain:latest-aws`
```
make docker_aws
```

Once that is completed, you just need to run the following command
which will create your container with a running instance of lightchain.  
```
export HOST_DATADIR=${HOME}/.lightchain_docker

docker run -v "${HOST_DATADIR}:/srv/lightchain" \
	-e "NETWORK=sirius" \
	-p 8545:8545 -p 26657:26657 -p 26656:26656 \
	-it lightchain:latest
```

Where `${HOST_DATADIR}` is the path in the host disk to
persist lightchain data. In case a volume is not used the node data
will be lost after the container is closed.

Alternatively to `sirius`, as NETWORK, you can also choose `mainnet` and `standalone`. 

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

### Running Tracer asserting blockchain internals, state, DB
[Read how to run the Tracer](TRACER.md)

### Running Unit Tests
```bash
go test ./...
```

## Credit & Licenses

- go-ethereum, [source](https://github.com/ethereum/go-ethereum),[license](https://github.com/ethereum/go-ethereum/#license)
- tendermint, [source](https://github.com/tendermint/tendermint),[license](https://github.com/tendermint/tendermint/blob/master/LICENSE)

## Bugs, Issues, Questions
If you find any bugs or simply have a question, please [write an issue](https://github.com/lightstreams-network/lightchain/issues) and we'll try and help as best we can.   

## Help
Need help? have questions? any feedback? Get in touch with us.
* [Discord](https://discordapp.com/invite/FXHZUKX)
* [Discuss Forum](https://discuss.lightstreams.network/c/dev)

## Contributors

- [Lukáš Lukáč](https://github.com/EnchanterIO), [Twitter](https://twitter.com/BlocksByLukas)
- [Gabriel Garrido](https://github.com/ggarri)
- [Michael Smolenski](https://github.com/mikesmo)
