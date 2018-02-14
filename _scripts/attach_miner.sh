#! /bin/bash

geth --exec 'miner.start()' attach ipc:$ETH_PATH/node$1/geth.ipc

