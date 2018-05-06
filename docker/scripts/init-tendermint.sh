#!/bin/bash

CONFIG_PATH="/tendermint/${HOSTNAME}/config"

#SEED_HOST="$(echo `hostname` | awk -F- '{print $1}')-0"
SEED_HOST="genesis-tendermint"

SEED_CONFIG_PATH="/tendermint/${SEED_HOST}/config"


#PUB_KEY=`jq '.pub_key' ${CONFIG_PATH}/priv_validator.json`
#PUBLIC_KEY=`jq '.pub_key.data' ${CONFIG_PATH}/priv_validator.json`
#grep --color ${PUBLIC_KEY} ${SEED_CONFIG_PATH}/genesis.json
#if [[ $? -eq 0 ]]; then
#	echo found
#else
#	echo not found
#	tmp=$(mktemp)
#	jq ".validators += [{pub_key: ${PUB_KEY}, power: 10, name: \"\"}]" ${SEED_CONFIG_PATH}/genesis.json > $tmp && mv $tmp ${SEED_CONFIG_PATH}/genesis.json
#fi
	
if [[ $HOSTNAME =~ ^.*-0$ ]] || [[ $HOSTNAME == 'genesis-tendermint' ]]; then
	echo true
	tendermint init --home=/tendermint/${HOSTNAME}
# v0.18.0
	#tendermint show_node_id --home=/tendermint/${HOSTNAME} > /tendermint/seed_id
# v0.16.0
        echo '' > /tendermint/seed_id
	cp ${SEED_CONFIG_PATH}/genesis.json /tendermint/genesis.json

else
	tendermint init --home=/tendermint/${HOSTNAME}
	cp /tendermint/genesis.json ${CONFIG_PATH}/genesis.json
	echo false
fi

$@
