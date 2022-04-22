#!/bin/bash
mkdir -p ./archives

source .env

docker-compose pull

docker save --output ./archives/core-command:${EDGEX_FOUNDRY_VERSION}.tar edgexfoundry/core-command:${EDGEX_FOUNDRY_VERSION} || exit 2

docker save --output ./archives/consul-${CONSUL_VERSION}.tar consul:${CONSUL_VERSION}

docker save --output ./archives/core-data-${EDGEX_FOUNDRY_VERSION}.tar edgexfoundry/core-data:${EDGEX_FOUNDRY_VERSION}

docker save --output ./archives/redis-${REDIS_VERSION}.tar redis:${REDIS_VERSION}

docker save --output ./archives/core-metadata-${EDGEX_FOUNDRY_VERSION}.tar edgexfoundry/core-metadata:${EDGEX_FOUNDRY_VERSION}

docker save --output ./archives/eclipse-mosquitto-${ECLIPSE_MOSQUITTO_VERSION}.tar eclipse-mosquitto:${ECLIPSE_MOSQUITTO_VERSION}

docker save --output ./archives/support-notifications-${EDGEX_FOUNDRY_VERSION}.tar edgexfoundry/support-notifications:${EDGEX_FOUNDRY_VERSION}

docker save --output ./archives/support-scheduler-${EDGEX_FOUNDRY_VERSION}.tar edgexfoundry/support-scheduler:${EDGEX_FOUNDRY_VERSION}

docker save --output ./archives/sys-mgmt-agent-${EDGEX_FOUNDRY_VERSION}.tar edgexfoundry/sys-mgmt-agent:${EDGEX_FOUNDRY_VERSION}

docker save --output ./archives/labs-air-edgex-app-service-metadata-${LABS_AIR_VERSION}.tar public.ecr.aws/tibcolabs/labs-air-edgex-app-service-metadata:${LABS_AIR_VERSION}

docker save --output ./archives/labs-air-edgex-device-generic-rest-${LABS_AIR_VERSION}.tar public.ecr.aws/tibcolabs/labs-air-edgex-device-generic-rest:${LABS_AIR_VERSION}

docker save --output ./archives/labs-air-edgex-device-generic-mqtt-${LABS_AIR_VERSION}.tar public.ecr.aws/tibcolabs/labs-air-edgex-device-generic-mqtt:${LABS_AIR_VERSION}