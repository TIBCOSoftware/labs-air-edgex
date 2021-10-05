// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_sound_simulator "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-sound-simulator"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-sound-simulator/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "edgex-sound-simulator"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_sound_simulator.Version, sd)
}
