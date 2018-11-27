#!/usr/bin/env bash

echo -e "Killing Lightchain"
ps -aux | grep "lightchain" | grep "node" | awk -F ' ' '{print $2}'| xargs kill -9
