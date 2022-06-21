#!/bin/bash

docker-compose -p edgex -f docker-compose-air-no-secty.yml -f docker-compose-air-no-secty-demo-generic.yml -f docker-compose-edgex-no-secty.yml down -v