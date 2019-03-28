#!/usr/bin/env bash

source $(dirname $0)/utils.sh

ROOT_PATH="$(cd "$(dirname "$0")" && pwd)/.."

DATA_DIR="${HOME}/.lightchain"
EXEC_BIN="./build/lightchain"
APPENDED_ARGS=""
NETWORK="sirius"

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
            NETWORK="standalone"
        ;;
        --mainnet) 
            NETWORK="mainnet" 
        ;;
        --sirius) 
            NETWORK="sirius" 
        ;;
        * )
            APPENDED_ARGS="${APPENDED_ARGS} $1"
    esac
    shift
done

DATA_DIR="${DATA_DIR}/${NETWORK}"
INIT_ARGS="--datadir=${DATA_DIR} --${NETWORK}"


pushd "$ROOT_PATH"

echo -e "Compiling latest version...."
if [ -n "${IS_DEBUG}" ]; then
	INIT_ARGS="${INIT_ARGS} --lvl=debug"
    run "make build-dev"
else
	INIT_ARGS="${INIT_ARGS} --lvl=info"
    run "make build"
fi

if [ -n "${CLEAN}" ]; then
	echo -e "You are about to wipe out ${DATA_DIR}"
    read -p "Are you sure? [N/y]" -n 1 -r
	echo    # (optional) move to a new line
	if [[ $REPLY =~ ^[Yy]$ ]]; then
	    echo -e "\t Restart environment"
	    echo "################################"
	    run "rm -rf ${DATA_DIR}"
	    INIT_ARGS="${INIT_ARGS} --trace"
		echo -e "################################ \n"
	else
		echo -e "Exiting"
		exit 1
	fi
fi

if [ -n "${IS_DEBUG}" ]; then
    EXEC_CMD="dlv --listen=:2345 --headless=true --api-version=2 exec ${EXEC_BIN} -- init ${INIT_ARGS}"
else
    EXEC_CMD="${EXEC_BIN} init ${INIT_ARGS}"
fi

run "$EXEC_CMD"

echo "Restoring ${NETWORK} private keys"
run "cp ./setup/${NETWORK}/database/keystore/* ${DATA_DIR}/database/keystore/"		

popd

echo -e "Execution completed"
exit 0
