#!/bin/bash

if [ -n "$1" ]; then
    case "$1" in
    arm64)
    export ARCH=-$1
    ;;
    amd64)
    export ARCH=""
    ;;
    esac
else
	export ARCH=""
fi

if [ -n "$2" ]; then
    EDGE_TYPE=-$2
else
    EDGE_TYPE="generic"
fi

docker-compose -p edgex -f docker-compose-air-no-secty.yml -f docker-compose-air-no-secty-demo-${EDGE_TYPE}.yml -f docker-compose-edgex-no-secty.yml down -v