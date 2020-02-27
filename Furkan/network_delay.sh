#!/bin/bash
if [ "$#" -ne 3 ]; then
    echo "Arguments: DURATION (in seconds, with no s), DELAY_TIME (in milliseconds), NODE_NAME"
    exit 1
fi
DURATION=$1
DELAY_TIME=$2
NODE_NAME=$3
sudo pumba netem --duration ${DURATION}s --interface eth0 --tc-image 'gaiadocker/iproute2' delay --time $DELAY_TIME $NODE_NAME
