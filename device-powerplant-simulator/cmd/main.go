// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_powerplant_simulator "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-powerplant-simulator"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-powerplant-simulator/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "edgex-powerplant-simulator"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_powerplant_simulator.Version, sd)
}
