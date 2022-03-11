@echo off
set COMPOSE_FILE=docker-compose.yml
echo "Using compose file: %COMPOSE_FILE%"

echo  " Starting.. %COMPOSE_FILE%"
docker-compose -f "%COMPOSE_FILE%" up

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

