// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_camera_onvif "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-camera-onvif"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-camera-onvif/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "edgex-device-camera-onvif"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_camera_onvif.Version, sd)
}
