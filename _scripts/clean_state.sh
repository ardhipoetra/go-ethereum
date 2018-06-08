#!/usr/bin/env bash

geth --datadir /home/ardhipoetra/git/eth/gik/storage//node0/ removedb
geth --datadir /home/ardhipoetra/git/eth/gik/storage//node1/ removedb

geth --datadir /home/ardhipoetra/git/eth/gik/storage/node0/ init /home/ardhipoetra/git/eth/gik/storage/head/genesis.json
geth --datadir /home/ardhipoetra/git/eth/gik/storage/node1/ init /home/ardhipoetra/git/eth/gik/storage/head/genesis.json