#! /bin/bash

geth --exec "$2" attach ipc:$ETH_PATH/node$1/geth.ipc

