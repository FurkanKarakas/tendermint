#!/bin/bash
if [ "$#" -ne 2 ]; then
    echo "Arguments: PORT_NR, TX_NAME"
    exit 1
fi
PORT_NR=$1
TX_NAME=$2
curl -s "localhost:${PORT_NR}/broadcast_tx_commit?tx=\"${TX_NAME}\""
