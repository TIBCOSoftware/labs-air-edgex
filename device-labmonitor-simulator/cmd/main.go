// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_labmonitor_simulator "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-labmonitor-simulator"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-labmonitor-simulator/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "edgex-labmonitor-simulator"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_labmonitor_simulator.Version, sd)
}
