#!/usr/bin/env bash

echo -e "Killing Lightchain"
ps -aux | grep "lightchain" | grep "dlv" | awk -F ' ' '{print $2}'| xargs kill -9
ps -aux | grep "lightchain" | awk -F ' ' '{print $2}'| xargs kill -9
