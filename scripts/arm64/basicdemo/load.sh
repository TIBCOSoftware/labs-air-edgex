#!/bin/bash

source .env

docker load --input ./archives/consul:${CONSUL_VERSION}.tar || exit 1

docker load --input ./archives/redis:${REDIS_VERSION}.tar || exit 1

docker load --input ./archives/core-metadata-arm64:${EDGEX_FOUNDRY_VERSION}.tar || exit 1

docker load --input ./archives/eclipse-mosquitto:${ECLIPSE_MOSQUITTO_VERSION}.tar || exit 1

docker load --input ./archives/core-command-arm64:${EDGEX_FOUNDRY_VERSION}.tar || exit 1

docker load --input ./archives/core-data-arm64:${EDGEX_FOUNDRY_VERSION}.tar || exit 1

docker load --input ./archives/labs-air-edgex-app-service-metadata-arm64:${LABS_AIR_VERSION}.tar || exit 1

docker load --input ./archives/labs-air-edgex-device-generic-rest-arm64:${LABS_AIR_VERSION}.tar || exit 1

docker load --input ./archives/labs-air-edgex-device-esp32-mqtt-arm64:${LABS_AIR_VERSION}.tar || exit 1

