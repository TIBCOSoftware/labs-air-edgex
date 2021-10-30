@echo off
set minikube_ip=$(minikube ip)
export AIR_MQTT_HOSTNAME=%minikube_ip%
call start.cmd
