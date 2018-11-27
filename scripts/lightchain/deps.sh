#!/usr/bin/env bash

ROOT_PATH="$(cd "$(dirname "$0")" && pwd)/.."

create_link() {
    VENDOR_RELATIVE_PATH="$1"
    CMD_RM="rm -rf ${ROOT_PATH}/vendor/${VENDOR_RELATIVE_PATH}"
    CMD_LK="ln -s  ${GOPATH}/src/${VENDOR_RELATIVE_PATH} ${ROOT_PATH}/vendor/${VENDOR_RELATIVE_PATH}"
    echo $CMD_RM
    eval $CMD_RM
    echo $CMD_LK
    eval $CMD_LK
}

echo -e "Hook up go-ethereum for go env sources"
create_link "github.com/ethereum/go-ethereum"
rm -rf "${ROOT_PATH}/vendor/github.com/ethereum/go-ethereum/vendor"

echo -e "Hook up go-ethereum for go env sources"
create_link "github.com/tendermint/tendermint"
rm -rf "${ROOT_PATH}/vendor/github.com/tendermint/tendermint/vendor"
