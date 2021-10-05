// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_camera_rtsp "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-camera-rtsp"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-camera-rtsp/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "edgex-device-camera-rtsp"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_camera_rtsp.Version, sd)
}
