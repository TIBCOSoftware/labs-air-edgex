// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	b64 "encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	sdkModel "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	device "github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
)

var energies = []string{
	"-7.881070663173226", "-7.681970643163226", "-7.591034687543226", "-7.618532746543298", "-7.607645600934267", "-7.596543890534876", "-7.675942040997654", "-7.643400926665487", "-7.705673456093245", "-7.730345218763976",
}

var distances = []string{
	"1.1", "1.2", "1.3", "1.2", "1.4", "1.0", "1.5", "1.4", "1.7", "1.4",
}

var currentts int64 = 0

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
	lc.Info(fmt.Sprintf("Initialize Device: %s", "SoundDevices"))
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

	lc.Info(fmt.Sprintf("Device publish interval: %d", interval))

	// Start Simulator
	// Set ticker to trigger data simulation
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {

		for {
			select {
			case <-ticker.C:
				lc.Info(fmt.Sprintf("Publishing simulated data"))
				publishSimulatedData()
			}
		}

	}()

	return nil
}

// HandleReadCommands executes a oommand
func (d *Driver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModel.CommandRequest) ([]*sdkModel.CommandValue, error) {
	d.Logger.Info(fmt.Sprintf("Device %s is handling read command", deviceName))

	var responses = make([]*sdkModel.CommandValue, len(reqs))
	var resTime = time.Now().UnixNano() / int64(time.Millisecond)
	var cv *sdkModel.CommandValue
	val := ""
	for i, req := range reqs {
		d.Logger.Debug(fmt.Sprintf("Request for resource: %s", req.DeviceResourceName))

		if req.DeviceResourceName == "Energy" {
			val = energies[rand.Intn(10)]
			d.Logger.Debug(fmt.Sprintf("Device %s is sending: %s", deviceName, val))
			cv, _ = sdkModel.NewCommandValueWithOrigin("Energy", common.ValueTypeString, val, resTime)
			responses[i] = cv
		} else if req.DeviceResourceName == "Distance" {
			val = distances[rand.Intn(10)]
			d.Logger.Debug(fmt.Sprintf("Device %s is sending: %s", deviceName, val))
			cv, _ = sdkModel.NewCommandValueWithOrigin("Distance", common.ValueTypeString, val, resTime)
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
	d.Logger.Info(fmt.Sprintf("A new Device is added: %s", deviceName))

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
	commandValues := make([]*sdkModel.CommandValue, 1)
	ind := GetcurrentSoundIndex()

	if currentts == 0 {
		currentts = time.Now().UnixNano() / int64(time.Millisecond)
	}

	for deviceName := range dynamicSimDevices {

		// Sound data
		// String data needed to be encoded to prevent issues in downstream
		// system parsing json messoge (i.e GraphBuilder to DGraph)
		val = b64.StdEncoding.EncodeToString([]byte(GetSoundData(ind)))
		cv, _ = sdkModel.NewCommandValueWithOrigin("microphone", common.ValueTypeString, val, currentts)
		commandValues[0] = cv

		asyncValues = &sdkModel.AsyncValues{
			DeviceName:    deviceName,
			CommandValues: commandValues,
		}
		driver.AsyncCh <- asyncValues

	}

	currentts = currentts + 60000
}
