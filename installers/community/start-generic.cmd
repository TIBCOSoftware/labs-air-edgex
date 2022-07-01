@echo off

echo %1

if "%1%"==""  set ARCH=""
if "%1%"=="amd64"  set ARCH=""
if "%1%"=="arm64"  set ARCH="-arm64"

echo "Using compose file: %COMPOSE_FILE%"

echo  " Starting Edge"


docker-compose -p edgex -f docker-compose-edgex-no-secty.yml  -f docker-compose-air-no-secty.yml -f docker-compose-air-no-secty-demo-generic.yml up


echo "\n"
echo -e " All services started. Edgex is ready "

