#!/bin/bash

extra_args=${1:-' --proxy_app=nilapp'}

tendermint node --consensus.create_empty_blocks=false --home=/tendermint/$HOSTNAME --p2p.seeds="${SEED_HOST}:46656" ${extra_args}
