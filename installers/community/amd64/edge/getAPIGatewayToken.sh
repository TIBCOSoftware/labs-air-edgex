
#!/bin/sh

CORE_EDGEX_REPOSITORY=edgexfoundry
CORE_EDGEX_VERSION=2.2.0

# Init key dir
GW_KEY_DIR=${GW_KEY_DIR:-${PWD}/keys}

echo ${GW_KEY_DIR}

# JWT File
JWT_FILE=/tmp/edgex/secrets/security-proxy-setup/kong-admin-jwt
JWT_VOLUME=/tmp/edgex/secrets/security-proxy-setup

ID=`cat ${GW_KEY_DIR}/gateway.id`

echo ${ID}

docker run --rm -it -e KONGURL_SERVER=edgex-kong -e "ID=${ID}" --network edgex_edgex-network --entrypoint "" -v ${GW_KEY_DIR}:/keys \
       ${CORE_EDGEX_REPOSITORY}/security-proxy-setup${ARCH}:${CORE_EDGEX_VERSION}${DEV} \
       /bin/sh -c '/edgex/secrets-config proxy jwt --algorithm ES256 --id ${ID} --private_key /keys/gateway.key'

