// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a simple example of a device service.
package main

import (
	device_opcua_go "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-opcua-go"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-opcua-go/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "edgex-device-opcua"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_opcua_go.Version, sd)
}
