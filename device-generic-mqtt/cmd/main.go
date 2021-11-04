// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_generic_mqtt "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-generic-mqtt"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-generic-mqtt/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "device-generic-mqtt"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_generic_mqtt.Version, sd)
}
