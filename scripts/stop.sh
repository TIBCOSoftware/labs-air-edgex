#!/bin/bash

network_type=${1:?}
os_type=${2:?}
arch_type=${3:?}
edge_type=${4:?}

# pushd ./${arch_type}/edge > /dev/null || exit 1
./stop-edge.sh ${arch_type} ${edge_type} || exit 2
# popd || exit 1