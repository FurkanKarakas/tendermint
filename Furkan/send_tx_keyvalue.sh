#!/bin/bash
if [ "$#" -ne 2 ]; then
    echo "Arguments: PORT_NR, TX_KEY, TX_VALUE"
    exit 1
fi
PORT_NR=$1
TX_KEY=$2
TX_VALUE=$3
curl -s "localhost:${PORT_NR}/broadcast_tx_commit?tx=\"${TX_KEY}=${TX_VALUE}\""
