#!/bin/bash


docker-compose -p edgex -f docker-compose-edgex-no-secty.yml  -f docker-compose-air-no-secty.yml -f docker-compose-air-no-secty-demo-ublox.yml up -d

echo $'\n'
echo -e "\033[0;32m All services started. Edge is ready\033[0m"

