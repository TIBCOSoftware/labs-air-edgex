// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_generic_mqtt "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-video-analytics-mqtt"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-video-analytics-mqtt/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/startup"
)

const (
	serviceName string = "device-video-analytics-mqtt"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_generic_mqtt.Version, sd)
}
