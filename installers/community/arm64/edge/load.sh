#!/bin/bash

# shellcheck source=/dev/null
source .env

docker load --input "./archives/arm64v8_consul:${CONSUL_VERSION}.tar" || exit 1

docker load --input "./archives/arm64v8_redis:${REDIS_VERSION}.tar" || exit 1

docker load --input "./archives/core-metadata-arm64:${EDGEX_FOUNDRY_VERSION}.tar" || exit 1

docker load --input "./archives/arm64v8_eclipse-mosquitto:${ECLIPSE_MOSQUITTO_VERSION}.tar" || exit 1

docker load --input "./archives/core-command-arm64:${EDGEX_FOUNDRY_VERSION}.tar" || exit 1

docker load --input "./archives/core-data-arm64:${EDGEX_FOUNDRY_VERSION}.tar" || exit 1

docker load --input "./archives/labs-air-edgex-app-service-metadata-arm64:${LABS_AIR_VERSION}.tar" || exit 1

docker load --input "./archives/labs-air-edgex-device-generic-rest-arm64:${LABS_AIR_VERSION}.tar" || exit 1

docker load --input "./archives/labs-air-edgex-device-esp32-mqtt-arm64:${LABS_AIR_VERSION}.tar" || exit 1

