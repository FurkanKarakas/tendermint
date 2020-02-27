#!/bin/bash
if [ "$#" -ne 3 ]; then
    echo "Arguments: DURATION (in seconds, with no s), DUPLICATE_PERCENT, NODE_NAME"
    exit 1
fi
DURATION=$1
DUPLICATE_PERCENT=$2
NODE_NAME=$3
sudo pumba netem --duration ${DURATION}s --interface eth0 --tc-image 'gaiadocker/iproute2' duplicate --percent $DUPLICATE_PERCENT $NODE_NAME
