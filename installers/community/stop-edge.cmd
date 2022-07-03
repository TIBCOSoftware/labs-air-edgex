@echo off

if "%1%"==""  set ARCH=""
if "%1%"=="amd64"  set ARCH=""
if "%1%"=="arm64"  set ARCH="-arm64"

set EDGE_TYPE=%2
if "%2%"==""  set EDGE_TYPE="generic"

echo  " Stopping Edge: %EDGE_TYPE%"

docker-compose -p edgex -f docker-compose-air-no-secty.yml -f docker-compose-air-no-secty-demo-%EDGE_TYPE%.yml -f docker-compose-edgex-no-secty.yml down -v