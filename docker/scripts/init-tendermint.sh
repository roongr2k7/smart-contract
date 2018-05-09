#!/bin/bash

CONFIG_PATH="/tendermint/${HOSTNAME}/config"

SEED_HOST=${3:-genesis-tendermint}

SEED_CONFIG_PATH="/tendermint/${SEED_HOST}/config"


if [[ $HOSTNAME =~ ^.*-0$ ]] || [[ $HOSTNAME == 'genesis-tendermint' ]]; then
	tendermint init --home=/tendermint/${HOSTNAME}
        echo '' > /tendermint/seed_id
else
	tendermint init --home=/tendermint/${HOSTNAME}
	curl ${SEED_HOST}:46657/genesis | jq ".result.genesis" > ${CONFIG_PATH}/genesis.json
        cat ${CONFIG_PATH}/genesis.json
fi

export SEED_HOST

$@
