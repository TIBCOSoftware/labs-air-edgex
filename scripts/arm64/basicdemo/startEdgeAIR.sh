#!/bin/sh

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

run_service () {
    echo "\n"
    echo -e "\033[0;32mStarting.. $1\033[0m"
    docker-compose -f "$COMPOSE_FILE" up -d $1
	
    if [ "$1" = "config-seed" ]
    then
         while [ -z "$(curl -s http://localhost:8500/v1/kv/config/device-virtual\;docker/app.open.msg)" ]
         do
               sleep 1
         done
         echo "$1 has been completely started !"
         return
    fi

    if [ -z "$2" ]
    then
         sleep 1
         echo "$1 has been completely started !"
         return
    fi
    
    if [ -n "$2" ]
    then
        sleep $2
        echo "$1 has been completely started !"
        return
    fi
}




run_service service-export-mqtt 5

# run_service service-export-kafka 48505

run_service service-flogo-rules

# run_service device-mqtt

# run_service device-particle-mqtt

# run_service device-stm32l475-mqtt

# run_service device-enersys-mqtt

# run_service device-jabil-rest

run_service device-vehicle-simulator

run_service device-siemens-simulator

# sssssrun_service device-smartwires-simulator

# run_service device-siemens-rest

# run_service device-virtual

# run_service device-random

echo "\n"
echo -e "\033[0;32m All services started. Edgex is ready\033[0m"

