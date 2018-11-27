#!/usr/bin/env bash

echo -e "Killing Tendermint nodes"
ps -aux | grep "tendermint" | grep "node" | awk -F ' ' '{print $2}'| xargs kill -9
