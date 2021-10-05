// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_depthai_mqtt "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-depthai-mqtt"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-depthai-mqtt/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "device-depthai-mqtt"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_depthai_mqtt.Version, sd)
}
