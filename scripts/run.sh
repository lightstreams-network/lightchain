#!/usr/bin/env bash

source $(dirname $0)/utils.sh

ROOT_PATH="$(cd "$(dirname "$0")" && pwd)/.."

DATA_DIR="${HOME}/.lightchain"
EXEC_BIN="./build/lightchain"
APPENDED_ARGS=""

while [ "$1" != "" ]; do
    case $1 in
        --datadir) 
            shift
            DATA_DIR=$1
        ;;
        --debug) 
            IS_DEBUG=1 
        ;;
        --clean) 
            CLEAN=1 
        ;;
        --standalone) 
            STANDALONE_NET=1 
        ;;
        * )
            APPENDED_ARGS="${APPENDED_ARGS} $1"
    esac
    shift
done

if [ -n "${STANDALONE_NET}" ]; then
	DATA_DIR="${DATA_DIR}/standalone"
	INIT_ARGS="--datadir=${DATA_DIR} --standalone"
else
	DATA_DIR="${DATA_DIR}/sirius"
	INIT_ARGS="--datadir=${DATA_DIR} --sirius"
fi

RUN_ARGS="--datadir=${DATA_DIR}"
RUN_ARGS="${RUN_ARGS} --rpc --rpcaddr=0.0.0.0 --rpcport=8545 --rpcapi eth,net,web3,personal,admin"
RUN_ARGS="${RUN_ARGS} --tmt_rpc_port=26657 --tmt_proxy_port=26658 --tmt_p2p_port=26656"
RUN_ARGS="${RUN_ARGS} --prometheus"

pushd "$ROOT_PATH"

echo -e "Compiling latest version...."
if [ -n "${IS_DEBUG}" ]; then
	RUN_ARGS="${RUN_ARGS} --lvl=debug"
    run "make build-dev"
else
	RUN_ARGS="${RUN_ARGS} --lvl=info"
    run "make build"
fi

if [ -n "${CLEAN}" ]; then
	echo -e "You are about to wipe out ${DATA_DIR} "
    read -p "Are you sure? [N/y]" -n 1 -r
	echo    # (optional) move to a new line
	if [[ $REPLY =~ ^[Yy]$ ]]; then
	    echo -e "\t Restart environment"
	    echo "################################"
	    run "rm -rf ${DATA_DIR}"
		run "$EXEC_BIN init ${INIT_ARGS}"
		if [ -z "${STANDALONE_NET}" ]; then
			echo "Restoring sirius private keys"
			run "cp ./setup/sirius/database/keystore/* ${DATA_DIR}/database/keystore/"		
		fi
		echo -e "################################ \n"
	else
		echo -e "Exiting"
		exit 1
	fi
fi

if [ -n "${IS_DEBUG}" ]; then
    EXEC_CMD="dlv --listen=:2345 --headless=true --api-version=2 exec ${EXEC_BIN} -- run ${RUN_ARGS}"
else
    EXEC_CMD="${EXEC_BIN} run ${RUN_ARGS}"
fi

run "$EXEC_CMD"

popd

echo -e "Execution completed"
exit 0
