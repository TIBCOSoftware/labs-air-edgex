#!/bin/bash

echo $1 # arch

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

docker-compose -p edgex -f docker-compose-edgex-no-secty.yml  -f docker-compose-air-no-secty.yml -f docker-compose-air-no-secty-demo-sound.yml up -d

echo $'\n'
echo -e "\033[0;32m All services started. Edge is ready\033[0m"

