// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_dai_camera_mask_mqtt "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-dai-camera-mask-mqtt"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-dai-camera-mask-mqtt/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "device-dai-camera-mask-mqtt"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_dai_camera_mask_mqtt.Version, sd)
}
