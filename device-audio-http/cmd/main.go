// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	device_audio_http "github.com/TIBCOSoftware/labs-air/edgexfoundry/device-audio-http"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/device-audio-http/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "edgex-audio-http"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_audio_http.Version, sd)
}
