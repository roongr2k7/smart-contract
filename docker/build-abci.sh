#!/bin/bash

# TODO: change from personal docker registry to ndid
docker image build -f Dockerfile-abci -t roongr2k7/abci:0.1.0 .
