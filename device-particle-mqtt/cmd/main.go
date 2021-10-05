// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_particle_mqtt "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-particle-mqtt"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-particle-mqtt/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "edgex-particle-mqtt"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_particle_mqtt.Version, sd)
}
