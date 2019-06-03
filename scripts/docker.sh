#!/usr/bin/env bash

set -e

LC_DATADIR="/srv/lightchain"

if [ -z "${NETWORK}" ]; then
	echo "Using default network sirius"
	NETWORK="sirius"
fi

if [ ! -d ${LC_DATADIR} ]; then
	echo "Create lightchain datadir ${LC_DATADIR}"
	mkdir -p ${LC_DATADIR}
fi

if [ ! -d "${LC_DATADIR}/database" ]; then
	echo "Initialize lightchain node in ${LC_DATADIR}"
	lightchain init --datadir=${LC_DATADIR} --${NETWORK}
fi

lightchain run --datadir=${LC_DATADIR} --rpc --rpcaddr=0.0.0.0 --rpcport=8545 --rpcapi=eth,net,web3,personal,debug
