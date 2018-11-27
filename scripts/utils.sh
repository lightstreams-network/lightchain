#!/usr/bin/env bash

function run {
    CMD="$1"
    echo -e "${CMD}"
    eval "${CMD} || exit 1"
}
