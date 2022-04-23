#!/bin/bash

network_type=${1:?}
os_type=${2:?}
arch_type=${3:?}

load_offline() {
    if [[ "${arch_type}" == "amd64" ]]; then
        pushd ./linux/basicdemo > /dev/null || exit 1
        ./load.sh || exit 2
        popd || exit 1
    elif [[ "${arch_type}" == "arm64" ]]; then
        pushd ./arm64/basicdemo > /dev/null || exit 1
        ./load.sh || exit 2
        popd || exit 1
    fi
}

if [[ "${network_type}" == "offline" ]]; then
      load_offline || exit 1
fi

start(){
    if [[ "${arch_type}" == "amd64" ]]; then
        pushd ./linux/basicdemo > /dev/null || exit 1
        ./start.sh || exit 2
        popd || exit 1
    elif [[ "${arch_type}" == "arm64" ]]; then
        pushd ./arm64/basicdemo > /dev/null || exit 1
        ./start.sh || exit 2
        popd || exit 1
    fi
}

