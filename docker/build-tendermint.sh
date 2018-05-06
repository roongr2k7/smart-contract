#!/bin/bash

# TODO: change from personal docker registry to ndid
docker image build -f Dockerfile-Tendermint -t roongr2k7/tendermint:0.16.0 .
