// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_dai_camera_people_mqtt "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-dai-camera-people-mqtt"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-dai-camera-people-mqtt/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "device-dai-camera-people-mqtt"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_dai_camera_people_mqtt.Version, sd)
}
