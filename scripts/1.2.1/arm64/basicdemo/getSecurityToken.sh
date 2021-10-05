#!/bin/sh

docker-compose run --rm --entrypoint /edgex/security-proxy-setup edgex-proxy --init=false --useradd=tibuser --group=admin