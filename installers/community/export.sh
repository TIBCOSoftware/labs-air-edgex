#!/bin/bash
mkdir -p ./archives

# shellcheck source=/dev/null
source .env

if [ -n "$1" ]; then
    case "$1" in
    arm64)
    export ARCH=-$1
    export ARCH_PF="arm64v8/"
    ;;
    amd64)
    export ARCH=""
    export ARCH_PF=""
    ;;
    esac
else
	export ARCH=""
    export ARCH_PF=""
fi


docker pull "edgexfoundry/core-command${ARCH}:${EDGEX_FOUNDRY_VERSION}" || exit 1
docker save --output "./archives/core-command-${EDGEX_FOUNDRY_VERSION}.tar" "edgexfoundry/core-command${ARCH}:${EDGEX_FOUNDRY_VERSION}" || exit 1

docker pull "${ARCH_PF}consul:${CONSUL_VERSION}" || exit 1
docker save --output "./archives/consul-${CONSUL_VERSION}.tar" "${ARCH_PF}consul:${CONSUL_VERSION}" || exit 1

docker pull "edgexfoundry/core-data${ARCH}:${EDGEX_FOUNDRY_VERSION}" || exit 1
docker save --output "./archives/core-data-${EDGEX_FOUNDRY_VERSION}.tar" "edgexfoundry/core-data${ARCH}:${EDGEX_FOUNDRY_VERSION}" || exit 1

docker pull "${ARCH_PF}redis:${REDIS_VERSION}" || exit 1
docker save --output "./archives/redis-${REDIS_VERSION}.tar" "${ARCH_PF}redis:${REDIS_VERSION}" || exit 1

docker pull "edgexfoundry/core-metadata${ARCH}:${EDGEX_FOUNDRY_VERSION}" || exit 1
docker save --output "./archives/core-metadata-${EDGEX_FOUNDRY_VERSION}.tar" "edgexfoundry/core-metadata${ARCH}:${EDGEX_FOUNDRY_VERSION}" || exit 1

docker pull "${ARCH_PF}eclipse-mosquitto:${ECLIPSE_MOSQUITTO_VERSION}" || exit 1
docker save --output "./archives/eclipse-mosquitto-${ECLIPSE_MOSQUITTO_VERSION}.tar" "${ARCH_PF}eclipse-mosquitto:${ECLIPSE_MOSQUITTO_VERSION}" || exit 1

docker pull "edgexfoundry/support-notifications${ARCH}:${EDGEX_FOUNDRY_VERSION}" || exit 1
docker save --output "./archives/support-notifications-${EDGEX_FOUNDRY_VERSION}.tar" "edgexfoundry/support-notifications${ARCH}:${EDGEX_FOUNDRY_VERSION}" || exit 1

docker pull "public.ecr.aws/tibcolabs/labs-air-edgex-app-service-metadata${ARCH}:${LABS_AIR_VERSION}" || exit 1
docker save --output "./archives/labs-air-edgex-app-service-metadata-${LABS_AIR_VERSION}.tar" "public.ecr.aws/tibcolabs/labs-air-edgex-app-service-metadata${ARCH}:${LABS_AIR_VERSION}" || exit 1

docker pull "public.ecr.aws/tibcolabs/labs-air-edgex-device-generic-rest${ARCH}:${LABS_AIR_VERSION}" || exit 1
docker save --output "./archives/labs-air-edgex-device-generic-rest-${LABS_AIR_VERSION}.tar" "public.ecr.aws/tibcolabs/labs-air-edgex-device-generic-rest${ARCH}:${LABS_AIR_VERSION}" || exit 1

docker pull "public.ecr.aws/tibcolabs/labs-air-edgex-device-generic-mqtt${ARCH}:${LABS_AIR_VERSION}" || exit 1
docker save --output "./archives/labs-air-edgex-device-generic-mqtt-${LABS_AIR_VERSION}.tar" "public.ecr.aws/tibcolabs/labs-air-edgex-device-generic-mqtt${ARCH}:${LABS_AIR_VERSION}" || exit 1

docker pull "public.ecr.aws/tibcolabs/labs-air-edgex-device-jetmax-mqtt${ARCH}:${LABS_AIR_VERSION}" || exit 1
docker save --output "./archives/labs-air-edgex-device-jetmax-mqtt-${LABS_AIR_VERSION}.tar" "public.ecr.aws/tibcolabs/labs-air-edgex-device-jetmax-mqtt${ARCH}:${LABS_AIR_VERSION}" || exit 1

docker pull "public.ecr.aws/tibcolabs/labs-air-edgex-device-sound-simulator${ARCH}:${LABS_AIR_VERSION}" || exit 1
docker save --output "./archives/labs-air-edgex-device-sound-simulator-${LABS_AIR_VERSION}.tar" "public.ecr.aws/tibcolabs/labs-air-edgex-device-sound-simulator${ARCH}:${LABS_AIR_VERSION}" || exit 1

docker pull "public.ecr.aws/tibcolabs/labs-air-edgex-device-esp32-mqtt${ARCH}:${LABS_AIR_VERSION}" || exit 1
docker save --output "./archives/labs-air-edgex-device-esp32-mqtt-${LABS_AIR_VERSION}.tar" "public.ecr.aws/tibcolabs/labs-air-edgex-device-esp32-mqtt${ARCH}:${LABS_AIR_VERSION}" || exit 1

docker pull "public.ecr.aws/tibcolabs/labs-air-edgex-device-ublox-rest${ARCH}:${LABS_AIR_VERSION}" || exit 1
docker save --output "./archives/labs-air-edgex-device-ublox-rest-${LABS_AIR_VERSION}.tar" "public.ecr.aws/tibcolabs/labs-air-edgex-device-ublox-rest${ARCH}:${LABS_AIR_VERSION}" || exit 1
