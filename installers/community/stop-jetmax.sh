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

docker-compose -p edgex -f docker-compose-air-no-secty.yml -f docker-compose-air-no-secty-demo-jetmax.yml -f docker-compose-edgex-no-secty.yml down -v