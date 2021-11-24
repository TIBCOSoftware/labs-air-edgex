@echo off
set COMPOSE_FILE=docker-compose.yml
echo "Using compose file: %COMPOSE_FILE%"



set EDGEX_TOKEN=""

echo  " Starting.. %COMPOSE_FILE%"
docker-compose -f "%COMPOSE_FILE%" up -d security-secrets-setup

docker-compose -f "%COMPOSE_FILE%" up -d consul

docker-compose -f "%COMPOSE_FILE%" up -d database

docker-compose -f "%COMPOSE_FILE%" up -d notifications

docker-compose -f "%COMPOSE_FILE%" up -d metadata

docker-compose -f "%COMPOSE_FILE%" up -d vault-worker

docker-compose -f "%COMPOSE_FILE%" up -d security-bootstrap-database

docker-compose -f "%COMPOSE_FILE%" up -d mqtt-broker

docker-compose -f "%COMPOSE_FILE%" up -d vault

docker-compose -f "%COMPOSE_FILE%" up -d kong-db

docker-compose -f "%COMPOSE_FILE%" up -d kong-migration
timeout 2

docker-compose -f "%COMPOSE_FILE%" up -d kong
timeout 5

docker-compose -f "%COMPOSE_FILE%" up -d edgex-proxy

docker-compose -f "%COMPOSE_FILE%" up -d scheduler

docker-compose -f "%COMPOSE_FILE%" up -d system

docker-compose -f "%COMPOSE_FILE%" up -d data

docker-compose -f "%COMPOSE_FILE%" up -d command


echo "Replacing token!!"
docker-compose run --rm --entrypoint /edgex/security-proxy-setup edgex-proxy --init=false --useradd=tibuser --group=admin > secToken.txt
findstr /B /C:"the access token for user tibuser is" secToken.txt > setTokenLine.txt
FOR /F "tokens=8 delims= " %%G IN (setTokenLine.txt) DO (set token=%%G)
set EDGEX_TOKEN=%token:~0,-1%


docker-compose -f "%COMPOSE_FILE%" up -d service-metadata
timeout 5

::KEEP DISSABLED: run_service device-siemens-simulator

docker-compose -f "%COMPOSE_FILE%" up -d device-generic-mqtt

docker-compose -f "%COMPOSE_FILE%" up -d device-generic-rest

echo "\n"
echo  "Installing cors plugin"
curl -X POST http://localhost:8001/plugins/ ^
--data "name=cors"  ^
--data "config.origins=*" ^
--data "config.methods=GET" ^
--data "config.methods=POST" ^
--data "config.methods=OPTIONS" ^
--data "config.headers=Accept" ^
--data "config.headers=Accept-Version" ^
--data "config.headers=Content-Length" ^
--data "config.headers=Content-MD5" ^
--data "config.headers=Content-Type" ^
--data "config.headers=Date" ^
--data "config.headers=X-Auth-Token" ^
--data "config.headers=Authorization" ^
--data "config.exposed_headers=X-Auth-Token" ^
--data "config.credentials=true" ^
--data "config.max_age=3600"


echo "\n"
echo -e " All services started. Edgex is ready "

