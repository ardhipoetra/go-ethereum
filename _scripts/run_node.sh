#!/usr/bin/env bash
# how to run : run_node <node_id> <logging directory> <[mine/nomine]> <DMCK ipc directory>
echo "Running the run_node.sh"
curdate=`date +"%Y%m%dT%H%M"`

node_id=$1
log_dir=$ETH_PATH/node$node_id/
ipc_dir=""

opt=" --selfId=$node_id"

if [ -n "$2" ]; then
  if [ "$2" == "mine" ]; then
    opt="$opt --mine --minerthreads=1"
  else
    log_dir=$2
  fi
fi

if [ -n "$3" ]; then
  if [ "$3" == "mine" ]; then
    opt="$opt --mine --minerthreads=1"
  fi
fi

if [ -n "$4" ]; then
  opt="$opt --dmckipcdir=$4"
fi

datadir=$ETH_PATH"node"$node_id

echo "Logging directory : $log_dir"
echo "Option : $opt"
echo "Data dir : $datadir"

set -o xtrace
geth $opt --networkid=1719 --port=3210$node_id --nodiscover\
    --verbosity=6 --vmodule="p2p=2"\
 --datadir=$datadir --identity="node$node_id"  --rpc --rpcport=990$node_id 2> "$log_dir/$curdate.log"

# tailf $ETH_PATH/node$node_id/$curdate.log
