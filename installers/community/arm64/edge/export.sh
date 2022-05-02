#!/bin/bash
mkdir -p ./archives

source .env

docker-compose pull || exit 1

docker save --output ./archives/consul:${CONSUL_VERSION}.tar consul:${CONSUL_VERSION} || exit 1

docker save --output ./archives/redis:${REDIS_VERSION}.tar redis:${REDIS_VERSION} || exit 1

docker save --output ./archives/core-metadata-arm64:${EDGEX_FOUNDRY_VERSION}.tar edgexfoundry/core-metadata-arm64:${EDGEX_FOUNDRY_VERSION} || exit 1

docker save --output ./archives/eclipse-mosquitto:${ECLIPSE_MOSQUITTO_VERSION}.tar eclipse-mosquitto:${ECLIPSE_MOSQUITTO_VERSION} || exit 1

docker save --output ./archives/core-command-arm64:${EDGEX_FOUNDRY_VERSION}.tar edgexfoundry/core-command-arm64:${EDGEX_FOUNDRY_VERSION} || exit 1

docker save --output ./archives/core-data-arm64:${EDGEX_FOUNDRY_VERSION}.tar edgexfoundry/core-data-arm64:${EDGEX_FOUNDRY_VERSION} || exit 1

docker save --output ./archives/labs-air-edgex-app-service-metadata-arm64:${LABS_AIR_VERSION}.tar public.ecr.aws/tibcolabs/labs-air-edgex-app-service-metadata-arm64:${LABS_AIR_VERSION} || exit 1

docker save --output ./archives/labs-air-edgex-device-generic-rest-arm64:${LABS_AIR_VERSION}.tar public.ecr.aws/tibcolabs/labs-air-edgex-device-generic-rest-arm64:${LABS_AIR_VERSION} || exit 1

docker save --output ./archives/labs-air-edgex-device-esp32-mqtt-arm64:${LABS_AIR_VERSION}.tar public.ecr.aws/tibcolabs/labs-air-edgex-device-esp32-mqtt-arm64:${LABS_AIR_VERSION} || exit 1








