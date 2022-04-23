#!/bin/bash

source .env

docker load --input ./archives/core-command:${EDGEX_FOUNDRY_VERSION}.tar || exit 1

docker load --input ./archives/consul-${CONSUL_VERSION}.tar || exit 1

docker load --input ./archives/core-data-${EDGEX_FOUNDRY_VERSION}.tar || exit 1

docker load --input ./archives/redis-${REDIS_VERSION}.tar || exit 1

docker load --input ./archives/core-metadata-${EDGEX_FOUNDRY_VERSION}.tar || exit 1

docker load --input ./archives/eclipse-mosquitto-${ECLIPSE_MOSQUITTO_VERSION}.tar || exit 1

docker load --input ./archives/support-notifications-${EDGEX_FOUNDRY_VERSION}.tar || exit 1

docker load --input ./archives/support-scheduler-${EDGEX_FOUNDRY_VERSION}.tar || exit 1

docker load --input ./archives/sys-mgmt-agent-${EDGEX_FOUNDRY_VERSION}.tar || exit 1

docker load --input ./archives/labs-air-edgex-app-service-metadata-${LABS_AIR_VERSION}.tar || exit 1

docker load --input ./archives/labs-air-edgex-device-generic-rest-${LABS_AIR_VERSION}.tar || exit 1

docker load --input ./archives/labs-air-edgex-device-generic-mqtt-${LABS_AIR_VERSION}.tar || exit 1