#!/usr/bin/env bash

source $(dirname $0)/utils.sh

ROOT_PATH="$(cd "$(dirname "$0")" && pwd)/.."

DATA_DIR="${HOME}/.lightchain"
EXEC_BIN="./build/lightchain"
APPENDED_ARGS=""
NETWORK="standalone"
PUBKEY=""
ADDRESS=""

ACTION="$1"
shift

if [ "${ACTION}" != "validatorset-deploy" ] && [ "${ACTION}" != "validatorset-add" ] && [ "${ACTION}" != "validatorset-remove" ] && [ "${ACTION}" != "validatorset-list" ]; then
	echo "Invalid action value '${ACTION}'. Valid values: [validatorset-deploy|validatorset-add|validatorset-remove|validatorset-list]"
	exit 1
fi

while [ "$1" != "" ]; do
    case $1 in
        --datadir) 
            shift
            DATA_DIR=$1
        ;;
        --debug) 
            IS_DEBUG=1 
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
        --pubkey) 
            shift
            PUBKEY=$1 
        ;;
        --address) 
            shift
            ADDRESS=$1 
        ;;
        --owner) 
            shift
            OWNER=$1 
        ;;
        --password) 
            shift
            PASSWORD=$1 
        ;;
        * )
            APPENDED_ARGS="${APPENDED_ARGS} $1"
    esac
    shift
done

if [ "${NETWORK}" = "standalone" ]; then
	OWNER="0xc916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e"
	PASSWORD="WelcomeToSirius"
fi

if [ -z "${OWNER}" ]; then
	echo "Missing argument --owner"
	exit 1
fi

DATA_DIR="${DATA_DIR}/${NETWORK}"
INIT_ARGS="--datadir=${DATA_DIR} --${NETWORK}"

RUN_ARGS="--datadir=${DATA_DIR}"
RUN_ARGS="${RUN_ARGS} --owner ${OWNER}"

if [ "${ACTION}" == "validatorset-add" ] || [ "${ACTION}" == "validatorset-remove" ]; then
	if [ -z "${PUBKEY}" ]; then
		echo "Missing value for --pubkey"
		exit 1
	fi
	
	if [ -z "${ADDRESS}" ]; then
		echo "Missing value for --address"
		exit 1
	fi

	RUN_ARGS="${RUN_ARGS} --pubkey ${PUBKEY}"
	RUN_ARGS="${RUN_ARGS} --address ${ADDRESS}"
fi

if [ -n "${PASSWORD}" ]; then
	RUN_ARGS="${RUN_ARGS} --password ${PASSWORD}"
fi

pushd "$ROOT_PATH"

echo -e "Compiling latest version...."
if [ -n "${IS_DEBUG}" ]; then
	RUN_ARGS="${RUN_ARGS} --lvl=debug"
    run "make build-dev"
else
	RUN_ARGS="${RUN_ARGS} --lvl=info"
    run "make build"
fi


if [ -n "${IS_DEBUG}" ]; then
    EXEC_CMD="dlv --listen=:2345 --headless=true --api-version=2 exec ${EXEC_BIN} -- governance ${ACTION} ${RUN_ARGS}"
else
    EXEC_CMD="${EXEC_BIN} governance ${ACTION} ${RUN_ARGS}"
fi

run "$EXEC_CMD"

popd

echo -e "Execution completed"
exit 0
