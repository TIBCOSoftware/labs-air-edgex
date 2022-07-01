
#!/bin/sh

CORE_EDGEX_REPOSITORY=edgexfoundry
CORE_EDGEX_VERSION=2.2.0

echo ${PWD}

# Init key dir
GW_KEY_DIR=${GW_KEY_DIR:-${PWD}/keys}

echo ${GW_KEY_DIR}

# JWT File
JWT_FILE=/tmp/edgex/secrets/security-proxy-setup/kong-admin-jwt
JWT_VOLUME=/tmp/edgex/secrets/security-proxy-setup

ID=`cat ${GW_KEY_DIR}/gateway.id`

echo ${ID}

docker run --rm -it -e KONGURL_SERVER=edgex-kong -e "ID=${ID}" -e "JWT_FILE=${JWT_FILE}" --network edgex_edgex-network --entrypoint "" -v ${GW_KEY_DIR}:/keys -v ${JWT_VOLUME}:${JWT_VOLUME} \
       ${CORE_EDGEX_REPOSITORY}/security-proxy-setup${ARCH}:${CORE_EDGEX_VERSION}${DEV} \
        /bin/sh -c 'JWT=`cat ${JWT_FILE}`; /edgex/secrets-config proxy adduser --token-type jwt --id ${ID} --algorithm ES256  --public_key /keys/gateway.pub \
              --user gateway --group gateway --jwt ${JWT} > /dev/null'
