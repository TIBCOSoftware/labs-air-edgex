#!/bin/bash




install_cors_plugin () {
    echo $'\n'
    echo  "Installing cors plugin"
    curl -X POST http://localhost:8001/plugins/ \
    --data "name=cors"  \
    --data "config.origins=*" \
    --data "config.methods=GET" \
    --data "config.methods=POST" \
    --data "config.methods=OPTIONS" \
    --data "config.headers=Accept" \
    --data "config.headers=Accept-Version" \
    --data "config.headers=Content-Length" \
    --data "config.headers=Content-MD5" \
    --data "config.headers=Content-Type" \
    --data "config.headers=Date" \
    --data "config.headers=X-Auth-Token" \
    --data "config.headers=Authorization" \
    --data "config.exposed_headers=X-Auth-Token" \
    --data "config.credentials=true" \
    --data "config.max_age=3600"
}

# docker-compose -p edgex -f ${COMPOSE_FILE} up -d

# docker-compose -p edgex -f docker-compose-edgex.yml -f docker-compose-air.yml -f docker-compose-airdemo-generic.yml up -d

docker-compose -p edgex -f docker-compose-edgex-no-secty.yml  -f docker-compose-air-no-secty.yml -f docker-compose-air-no-secty-demo-generic.yml up -d

# Only needed when edgex is running with security enabled
# install_cors_plugin


echo $'\n'
echo -e "\033[0;32m All services started. Edgex is ready\033[0m"

