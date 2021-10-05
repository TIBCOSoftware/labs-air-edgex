// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_dai_camera_palm_mqtt "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-dai-camera-palm-mqtt"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-dai-camera-palm-mqtt/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "device-dai-camera-palm-mqtt"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_dai_camera_palm_mqtt.Version, sd)
}
