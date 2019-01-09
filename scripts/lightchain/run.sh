#!/usr/bin/env bash

source $(dirname $0)/../utils.sh

ROOT_PATH="$(cd "$(dirname "$0")" && pwd)/../.."

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

INIT_ARGS=""
#INIT_ARGS="--homedir ${DATA_DIR}"

NODE_ARGS=""
#NODE_ARGS="--homedir ${DATA_DIR}"
NODE_ARGS="${NODE_ARGS} --tmt_rpc_port 25557 --tmt_proxy_port=25558 tmt_p2p_port=25556"
NODE_ARGS="${NODE_ARGS} --rpc --rpcaddr=0.0.0.0 --rpcport=8545"
NODE_ARGS="${NODE_ARGS} --rpcapi eth,net,web3,personal,admin"
#NODE_ARGS="${NODE_ARGS} --ws --wsaddr=0.0.0.0 --wsport=8546 --wsorigins='*'"
#NODE_ARGS="${NODE_ARGS} --abci_laddr=tcp://0.0.0.0:26658 --tendermint_addr=tcp://127.0.0.1:26657"

pushd "$ROOT_PATH"

if [ -n "${CLEAN}" ]; then
    echo "################################"
    echo -e "\t Restart environment"
    echo "################################"
    if [ -n "${HARD_MODE}" ]; then
		run "rm -rf ${DATA_DIR}"
	else
		run "${EXEC_BIN} ${INIT_ARGS} unsafe_reset_all"
	fi
	run "${EXEC_BIN} ${INIT_ARGS} init"
    echo -e "################################ \n"
fi

if [ -n "${IS_DEBUG}" ]; then
    EXEC_CMD="dlv --listen=:2345 --headless=true --api-version=2 exec ${EXEC_BIN} -- ${NODE_ARGS} node"
else
    EXEC_CMD="${EXEC_BIN} ${NODE_ARGS} node"
fi

run "$EXEC_CMD"

popd

echo -e "Execution completed"
exit 0
