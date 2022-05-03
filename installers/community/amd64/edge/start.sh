#!/bin/bash

# Copyright 2017 Konrad Zapalowicz <bergo.torino@gmail.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Start EdgeX Foundry services in right order, as described:
# https://wiki.edgexfoundry.org/display/FA/Get+EdgeX+Foundry+-+Users

COMPOSE_FILE=${1:-docker-compose.yml}
echo "Using compose file: $COMPOSE_FILE"


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

docker-compose -p edgex -f ${COMPOSE_FILE} up -d

# Only needed when edgex is running with security enabled
# install_cors_plugin


echo $'\n'
echo -e "\033[0;32m All services started. Edgex is ready\033[0m"

