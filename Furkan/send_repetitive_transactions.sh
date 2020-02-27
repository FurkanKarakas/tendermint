#!/bin/bash
if [ "$#" -ne 4 ]; then
    echo "Arguments: PORT_NR, TX_NAME, TX_SIZE, TIME_BETWEEN_TX"
    exit 1
fi
PORT_NR=$1
TX_NAME=$2
TX_SIZE=$3
TIME_BETWEEN_TX=$4
for i in $(seq $TX_SIZE)
do
    curl -s "localhost:${PORT_NR}/broadcast_tx_commit?tx=\"${TX_NAME}$i\""
    sleep $TIME_BETWEEN_TX
done
