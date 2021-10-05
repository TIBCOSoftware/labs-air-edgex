// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_hololens_mqtt "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-hololens-mqtt"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-hololens-mqtt/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "device-hololens-mqtt"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_hololens_mqtt.Version, sd)
}
