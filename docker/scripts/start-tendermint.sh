#!/bin/bash

#SEED_HOST="$(echo `hostname` | awk -F- '{print $1}')-0"
#tendermint node --consensus.create_empty_blocks=false --home=/tendermint/${HOSTNAME} --p2p.seeds="`cat /tendermint/seed_id`@${SEED_HOST}:46656" --proxy_app=nilapp

extra_args=${@:-' --proxy_app=nilapp'}
seed_id=`cat /tendermint/seed_id`

if [[ -z ${seed_id} ]]; then
	tendermint node --consensus.create_empty_blocks=false --home=/tendermint/$HOSTNAME --p2p.seeds="genesis-tendermint:46656" ${extra_args}
else
	tendermint node --consensus.create_empty_blocks=false --home=/tendermint/$HOSTNAME --p2p.seeds="${seed_id}@genesis-tendermint:46656" ${extra_args}
fi

