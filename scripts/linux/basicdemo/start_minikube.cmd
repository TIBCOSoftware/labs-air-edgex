@echo off
FOR /F "tokens=1 delims= " %%G IN ('minikube ip') DO (set minikube_ip=%%G)
set AIR_MQTT_HOSTNAME=%minikube_ip%
call start.cmd
