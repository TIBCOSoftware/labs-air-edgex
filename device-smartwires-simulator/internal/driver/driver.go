// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	sdkModel "github.com/edgexfoundry/device-sdk-go/pkg/models"
	device "github.com/edgexfoundry/device-sdk-go/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

var voltages = []string{
	"220.0", "221.0", "225.0", "218.0", "220.5", "219.0", "224.0", "218.0", "223.0", "220.0",
}

var amperages = []string{
	"30.0", "29.0", "28.0", "30.0", "31.0", "32.0", "27.0", "29.0", "30.0", "31.0",
}

var temperatures = []string{
	"90.0", "120.0", "140.0", "130.0", "135.0", "140.0", "125.0", "130.0", "135.0", "120.0",
}

var once sync.Once
var driver *Driver
var dynamicSimDevices map[string]bool

// Driver the driver sructure
type Driver struct {
	Logger           logger.LoggingClient
	AsyncCh          chan<- *sdkModel.AsyncValues
	CommandResponses sync.Map
}

// NewProtocolDriver creates a new driver
func NewProtocolDriver() sdkModel.ProtocolDriver {
	once.Do(func() {
		driver = new(Driver)
	})
	return driver
}

// DisconnectDevice disconnect
func (d *Driver) DisconnectDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	d.Logger.Warn("Driver's DisconnectDevice function didn't implement")
	return nil
}

// Initialize the device
func (d *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *sdkModel.AsyncValues, deviceCh chan<- []sdkModel.DiscoveredDevice) error {
	d.Logger = lc
	d.AsyncCh = asyncCh

	// Initialize local device cache
	dynamicSimDevices = make(map[string]bool)

	// Get Driver Config
	interval := 30
	configMap := device.DriverConfigs()
	simulationInterval, ok := configMap["SimulationInterval"]
	if ok {
		interval, _ = strconv.Atoi(simulationInterval)
	}

	// Start Simulator
	// Set ticker to trigger data simulation
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		// If you want to start simulation right away
		// for ; true; <-ticker.C {
		// 	publishSimulatedData()
		// }

		for {
			select {
			case <-ticker.C:
				publishSimulatedData()
			}
		}

	}()

	return nil
}

// HandleReadCommands executes a oommand
// Note:  This simulator assumes there is always only one request (Status). However,
// with each status, there are other resource values that are also published
func (d *Driver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModel.CommandRequest) ([]*sdkModel.CommandValue, error) {
	d.Logger.Debug(fmt.Sprintf("Device %s is handling read command", deviceName))

	var responses = make([]*sdkModel.CommandValue, len(reqs))
	var resTime = time.Now().UnixNano() / int64(time.Millisecond)
	var cv *sdkModel.CommandValue
	val := ""
	for i, req := range reqs {
		d.Logger.Debug(fmt.Sprintf("Request for resource: %s", req.DeviceResourceName))

		if req.DeviceResourceName == "Voltage" {
			val = voltages[rand.Intn(10)]
			d.Logger.Debug(fmt.Sprintf("Device %s is sending: %s", deviceName, val))
			cv = sdkModel.NewStringValue("Voltage", resTime, val)
			responses[i] = cv
		} else if req.DeviceResourceName == "Amperage" {
			val = amperages[rand.Intn(10)]
			d.Logger.Debug(fmt.Sprintf("Device %s is sending: %s", deviceName, val))
			cv = sdkModel.NewStringValue("Amperage", resTime, val)
			responses[i] = cv
		} else if req.DeviceResourceName == "Temperature" {
			val = temperatures[rand.Intn(10)]
			d.Logger.Debug(fmt.Sprintf("Device %s is sending: %s", deviceName, val))
			cv = sdkModel.NewStringValue("Temperature", resTime, val)
			responses[i] = cv
		}

	}

	d.Logger.Debug(fmt.Sprintf("Device %s is sending response", deviceName))
	return responses, nil
}

// HandleWriteCommands write command
func (d *Driver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModel.CommandRequest, params []*sdkModel.CommandValue) error {
	var err error

	return err
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (d *Driver) Stop(force bool) error {
	// Then Logging Client might not be initialized
	if d.Logger != nil {
		d.Logger.Debug(fmt.Sprintf("Driver.Stop called: force=%v", force))
	}
	return nil
}

// AddDevice is a callback function that is invoked
// when a new Device associated with this Device Service is added
func (d *Driver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.Logger.Debug(fmt.Sprintf("a new Device is added: %s", deviceName))

	dynamicSimDevices[deviceName] = true

	return nil
}

// UpdateDevice is a callback function that is invoked
// when a Device associated with this Device Service is updated
func (d *Driver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.Logger.Debug(fmt.Sprintf("Device %s is updated", deviceName))
	return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (d *Driver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	d.Logger.Debug(fmt.Sprintf("Device %s is removed", deviceName))

	delete(dynamicSimDevices, deviceName)

	return nil
}

func publishSimulatedData() {

	val := ""
	var cv *sdkModel.CommandValue
	var asyncValues *sdkModel.AsyncValues
	var resTime = time.Now().UnixNano() / int64(time.Millisecond)

	for deviceName := range dynamicSimDevices {

		// Publish Voltage
		val = voltages[rand.Intn(10)]
		cv = sdkModel.NewStringValue("Voltage", resTime, val)

		asyncValues = &sdkModel.AsyncValues{
			DeviceName:    deviceName,
			CommandValues: []*sdkModel.CommandValue{cv},
		}
		driver.AsyncCh <- asyncValues

		// Publish Amperage
		val = amperages[rand.Intn(10)]
		cv = sdkModel.NewStringValue("Amperage", resTime, val)
		asyncValues = &sdkModel.AsyncValues{
			DeviceName:    deviceName,
			CommandValues: []*sdkModel.CommandValue{cv},
		}
		driver.AsyncCh <- asyncValues

		// Publish Temperature
		val = temperatures[rand.Intn(10)]
		cv = sdkModel.NewStringValue("Temperature", resTime, val)
		asyncValues = &sdkModel.AsyncValues{
			DeviceName:    deviceName,
			CommandValues: []*sdkModel.CommandValue{cv},
		}
		driver.AsyncCh <- asyncValues
	}

}
