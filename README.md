# Lightchain

Official Lightstreams blockchain based on go-ethereum and Tendermint.

## Installation

### Lightchain

```
export GIT_TERMINAL_PROMPT=1
go get -u github.com/lightstreams-network/lightchain
cd ${GOPATH}/src/github.com/lightstreams-network/lightchain
make get_vendor_deps
make install
```

### Install Tendermint 0.27.0 (patched)

Running the following commands to download and install Tendermint
in your local environment.

```
mkdir -p $GOPATH/src/github.com/tendermint
cd $GOPATH/src/github.com/tendermint
git clone https://github.com/lightstreams-network/tendermint.git
cd tendermint

make get_tools
make get_vendor_deps
make install
```

> Optionally you can also following official [docs](https://tendermint.com/docs/introduction/install.html):
  
## Launching Node

### Step 1: Initialization

#### Lightchain

To initialise a new blockchain you need to run `init` command and choose a local path where blockchain files are going to be stored.

```
lightchain --datadir "${HOME}/.lightchain" init [GENESIS_JSON_PATH]
```

By default it is used the `genesis.json` but alternatively you can select a genesis as a second argument.

```
├── lightchain
│   └── chaindata
│       ├── 000001.log
│       ├── CURRENT
│       ├── LOCK
│       ├── LOG
│       └── MANIFEST-000000
├── genesis.json
└── keystore
    ├── UTC--2018-08-26T21-40-07.289727986Z--4eaaad8ea38d5ef75ebdeb3d1be59d56f86c4ca9
    └── UTC--2018-08-26T21-40-31.689362077Z--4f5adedca6d869e9f5f7dcf4b7a9dfa8231a095f
```

- `/chaindata/*`: Blockchain database files
- `genesis.json`: Genesis block used to initialize blockchain.
- `keystore/*`:  Accounts private keys, those accounts have been initialized along with the blockchain, their balance is defined as part of the `genesis` block.

###### Optional (Connect to `sirius`) 

In the case we want to get our local tendermint node connected to `sirius` network
we need to replace the files `genesis.json` and `config.tolm` by 
the ones you find at the path `./setup/tendermint`
```
cp ./setup/tendermint/genesis.json ${HOME}/.lightchain/tendermint/config/
cp ./setup/tendermint/config.tolm ${HOME}/.lightchain/tendermint/config/
``` 

### Step 2: Start-up Tendermint server

```
tendermint --home ${HOME}/.lightchain/tendermint --consensus.create_empty_blocks=false node
```

In case you are running your node to be connected to `sirius`
it will take short while till it is accepted and synchronized.  

### Step 3: Start-up Lightchain server

```
lightchain --datadir ${HOME}/.lightchain --rpc --rpcaddr=0.0.0.0 --ws --wsaddr=0.0.0.0 --rpcapi eth,net,web3,personal,admin node
```

## Documentation

To know more about how Lightchain works and how Tendermint is integrated
to perform the PoA, visit our repository [wiki](https://github.com/lightstreams-network/lightchain/wiki)

## Tests

[Read how to run tests.](truffle/Tests.md)
