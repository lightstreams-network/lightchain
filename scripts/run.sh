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
        --hard) 
            HARD_MODE=1 
        ;;
        * )
            APPENDED_ARGS="${APPENDED_ARGS} $1"
    esac
    shift
done

INIT_ARGS="--datadir=${DATA_DIR}"

RUN_ARGS="--datadir=${DATA_DIR} --lvl=dbug"

pushd "$ROOT_PATH"

if [ -n "${CLEAN}" ]; then
    echo "################################"
    echo -e "\t Restart environment"
    echo "################################"
    if [ -n "${HARD_MODE}" ]; then
        run "rm -rf ${DATA_DIR}"
    fi
	run "$EXEC_BIN init ${INIT_ARGS}"
    echo -e "################################ \n"
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
