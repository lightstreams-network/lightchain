# Lightchain

Official Lightstreams blockchain based on go-ethereum and Tendermint.

## Installation

### Install Lightchain

```
go get -u github.com/lightstreams-network/lightchain
cd ${GOPATH}/src/github.com/lightstreams-network/lightchain
make get_vendor_deps
make install
```

### Build Lightchain

```
make build
```

OR to compile with debug flags

```
make build-dev
```

### Install Tendermint 0.25.1-rc0

Following official [docs](https://tendermint.com/docs/introduction/install.html):

```
mkdir -p $GOPATH/src/github.com/tendermint
cd $GOPATH/src/github.com/tendermint
git clone https://github.com/tendermint/tendermint.git
cd tendermint
git reset --hard tags/v0.25.1-rc0

make get_tools
make get_vendor_deps

make install
```

### Build Tendermint

```
make build
```
  
## Usage

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

#### Tendermint

Tendermint requires to be initialized, to do that first we need to define where tendermint data is going to be allocated, for the sake of simplicity it is recommend to be stored in a sub-folder of chosen lightchain datadir.

```
tendermint --home ${HOME}/.lightchain/tendermint init
```

```
├── config
│   ├── config.toml
│   ├── genesis.json
│   ├── node_key.json
│   └── priv_validator.json
└── data
```
 To know more about the meaning of those files please visit [Tendermint official doc](https://tendermint.com/docs/).

### Step 2: Start-up Tendermint server

```
tendermint --home ${HOME}/.lightchain/tendermint --consensus.create_empty_blocks=false node
```

### Step 3: Start-up Lightchain server

```
lightchain --datadir ${HOME}/.lightchain --rpc --rpcaddr=0.0.0.0 --ws --wsaddr=0.0.0.0 --rpcapi eth,net,web3,personal,admin node
```