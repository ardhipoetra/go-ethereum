#!/bin/bash

node_id=$1
opt=''
if [ -n "$2" ]; then
  if [ "$2" == "mine" ]; then
    opt='--mine --minerthreads=1'
  fi
fi

curdate=`date +"%Y%m%dT%H%M"`

geth --networkid=1719 --port=3210$node_id --nodiscover\
    --verbosity=6 --vmodule="p2p=2"\
 --datadir=$ETH_PATH/node$node_id --identity="node$node1"  --rpc --rpcport=990$node_id #2> $ETH_PATH/node$node_id/$curdate.log

# tailf $ETH_PATH/node$node_id/$curdate.log


