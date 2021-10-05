// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_wafer_simulator "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-wafer-simulator"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-wafer-simulator/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "edgex-wafer-simulator"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_wafer_simulator.Version, sd)
}
