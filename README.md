# Lightchain

This is the official Lightstreams implementation of a proof-of-authority (PoA) blockchain. Lightchain is an ethereum-compatible blockchain which uses byzantine consensus to replace the original proof-of-work (PoW) from Ethereum. This is achieved by integrating [`Tendermint`](https://tendermint.com) for the consensus layer. 

## Pre-requirements

- Go >= 1.10
- deps

## Installation


To install `lightchain` in your system just run following commands:
```
export GIT_TERMINAL_PROMPT=1
go get -u github.com/lightstreams-network/lightchain
cd ${GOPATH}/src/github.com/lightstreams-network/lightchain
make get_vendor_deps
make install
```

To validate if you installation was done correctly, run the following command to obtain current installed version of Lightchain
```
lightchain version
```
  
## Create a Lightchain node

Lightstreams provides a testnet called `Sirius` by default new created nodes are being hooked to this network and automatically synchronized. Currently lightstreams team is actively working on improving the network performance and stability therefore some issues might still occur which force us to restore blockchain. To check the current state of `Sirius` network we provide an [blockchain explorer](https://explorer.lightstreams.io/home) 

We are working on launching the **Lightstream main network** which aims to provide a reliable and fast open blockchain. Stay tuned about our project's progress by reading our [blog](https://medium.com/lightstreams)
 
### Node initialization

To initialise a new blockchain you need to run `lightchain init` and  choose a local path where blockchain files are going to be stored.
```
lightchain init --datadir "${HOME}/.lightchain"
```

Once this is done, your blockchain is ready to be launche. The created node is setup to connect to the lightstreams' test network (`Sirius`). If you prefer to create an isolated node for testing proposes, you can run the following command instead:
```
lightchain init --datadir "${HOME}/.lightchain" --standalone
```

When creating a standalone node, a faucet account with funds is created by default: 
```
@TODO
```

### Node launch

To run a lightchain node you only need to run the following command:
```
lightchain run --datadir "${HOME}/.lightchain"
```

After running the command above (and only if you are not in `standalone` mode), it will start running the network synchronization which will take several minutes.


***Available flags***
```
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

## Applications

### Leth
Lightstreams implemented its first DApp running onto Lightchain blockchain
[**Leth**](https://docs.lightstreams.network/01.getting-started/). 

Leth is application which intends to wrap Ethereum blockchain + IPFS into a very simple interface which can be used either by [HTTP Restful API](https://docs.lightstreams.network/api-docs/) or by [Interactive Command line client] (https://docs.lightstreams.network/04.cli-docs/leth/) 

## Docker
In case you prefer to use Docker, follow the instructions below.

First you create a new docker image, which will be tagged as `lightchain:latest`
```
make docker
```

Once the above execution is completed you just need to run the following statement
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
 


## Documentation

***Project data structure***

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

***Wiki***
To know more about how Lightchain works and how Tendermint is integrated to perform the PoA, visit our repository [wiki](https://github.com/lightstreams-network/lightchain/wiki)

## Tests
[Read how to run tests](truffle/Tests.md)
