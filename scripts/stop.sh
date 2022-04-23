#!/bin/bash

docker-compose -p edgex down -v
network_type=${1:?}
os_type=${2:?}
arch_type=${3:?}

if [[ "${arch_type}" == "amd64" ]]; then
    pushd ./linux/basicdemo > /dev/null || exit 1
    ./stop.sh || exit 2
    popd || exit 1
elif [[ "${arch_type}" == "arm64" ]]; then
    pushd ./arm64/basicdemo > /dev/null || exit 1
    ./stop.sh || exit 2
    popd || exit 1
fi