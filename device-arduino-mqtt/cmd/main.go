// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_arduino_mqtt "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-arduino-mqtt"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-arduino-mqtt/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "device-arduino-mqtt"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_arduino_mqtt.Version, sd)
}
