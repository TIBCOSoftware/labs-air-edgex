#!/bin/bash

minikube_ip=$(minikube ip)
export AIR_MQTT_HOSTNAME=${MQTT_HOSTNAME:-$minikube_ip}
source ./start.sh
