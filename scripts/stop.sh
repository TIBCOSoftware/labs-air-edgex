#!/bin/bash

network_type=${1:?}
os_type=${2:?}
arch_type=${3:?}

# pushd ./${arch_type}/edge > /dev/null || exit 1
./stop-generic.sh ${arch_type} || exit 2
# popd || exit 1