// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	sdkModel "github.com/edgexfoundry/device-sdk-go/pkg/models"
	device "github.com/edgexfoundry/device-sdk-go/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

// var sampleRate int32 = 0
// var audioStreamerURL string = "http://localhost:8080/audio"
// var sampleDurationSecs int = 2

var once sync.Once
var driver *Driver
var dynamicSimDevices map[string]bool

// Driver the driver sructure
type Driver struct {
	Logger                 logger.LoggingClient
	AsyncCh                chan<- *sdkModel.AsyncValues
	DeviceCh               chan<- []sdkModel.DiscoveredDevice
	CommandResponses       sync.Map
	StreamerURL            string
	StreamerSampleRate     int
	StreamerSampleDuration int
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
	lc.Info(fmt.Sprintf("Initialize Device: %s", "AudioDevices"))
	d.Logger = lc
	d.AsyncCh = asyncCh
	d.DeviceCh = deviceCh

	// Get Driver Config
	configMap := device.DriverConfigs()
	d.StreamerURL = configMap["AudioStreamerURL"]
	sampleRate := configMap["StreamerSampleRate"]
	d.StreamerSampleRate, _ = strconv.Atoi(sampleRate)
	sampleDuration := configMap["StreamerSampleDurationSec"]
	d.StreamerSampleDuration, _ = strconv.Atoi(sampleDuration)

	lc.Info(fmt.Sprintf("Audio Streamer URL: %s", d.StreamerURL))

	// Initialize local device cache
	dynamicSimDevices = make(map[string]bool)

	// Start Simulator
	// Set ticker to trigger data simulation
	// ticker := time.NewTicker(time.Duration(interval) * time.Second)
	// go func() {

	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			lc.Info("Publishing simulated data")
	// 			publishAudioData()
	// 		}
	// 	}

	// }()

	return nil
}

// HandleReadCommands executes a oommand
func (d *Driver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModel.CommandRequest) ([]*sdkModel.CommandValue, error) {
	d.Logger.Info(fmt.Sprintf("Device %s is handling read command", deviceName))

	var responses = make([]*sdkModel.CommandValue, len(reqs))
	var resTime = time.Now().UnixNano() / int64(time.Millisecond)

	d.Logger.Info(fmt.Sprintf("Stream server address: %s", d.StreamerURL))

	var cv *sdkModel.CommandValue
	val := ""
	for i, req := range reqs {
		d.Logger.Debug(fmt.Sprintf("Request for resource: %s", req.DeviceResourceName))

		if req.DeviceResourceName == "microphone" {

			data, err := getAudioData(d.StreamerURL, d.StreamerSampleRate*d.StreamerSampleDuration)

			if err != nil {
				d.Logger.Error(err.Error())
				return responses, err
			}

			cv, _ = sdkModel.NewFloat32ArrayValue("microphone", resTime, data)

			d.Logger.Debug(fmt.Sprintf("Device %s is sending: %s", deviceName, val))

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

func (d *Driver) addrFromProtocols(protocols map[string]models.ProtocolProperties) (string, error) {

	if _, ok := protocols["HTTP"]; !ok {
		d.Logger.Error("No HTTP address found for device. Check configuration file.")
		return "", fmt.Errorf("no HTTP address in protocols map")
	}

	var addr string
	addr, ok := protocols["HTTP"]["Address"]
	if !ok {
		d.Logger.Error("No HTTP address found for device. Check configuration file.")
		return "", fmt.Errorf("no RTSP address in protocols map")
	}
	return addr, nil

}

func getAudioData(addr string, dataSize int) ([]float32, error) {

	resp, err := http.Get(addr)

	if err == nil {

		body, _ := ioutil.ReadAll(resp.Body)
		responseReader := bytes.NewReader(body)
		buffer := make([]float32, dataSize)
		binary.Read(responseReader, binary.BigEndian, &buffer)

		return buffer, nil

	} else {
		return nil, fmt.Errorf(err.Error())
	}

}

// func publishAudioData() {

// 	val := ""
// 	var cv *sdkModel.CommandValue
// 	var asyncValues *sdkModel.AsyncValues
// 	commandValues := make([]*sdkModel.CommandValue, 1)

// 	resp, err := http.Get("http://localhost:8080/audio")

// 	if err == nil {

// 		body, _ := ioutil.ReadAll(resp.Body)

// 		currentts := time.Now().UnixNano() / int64(time.Millisecond)

// 		for deviceName := range dynamicSimDevices {

// 			// Sound data
// 			// String data needed to be encoded to prevent issues in downstream
// 			// system parsing json messoge (i.e GraphBuilder to DGraph)

// 			val = b64.StdEncoding.EncodeToString(body)
// 			cv = sdkModel.NewStringValue("microphone", currentts, val)
// 			commandValues[0] = cv

// 			asyncValues = &sdkModel.AsyncValues{
// 				DeviceName:    deviceName,
// 				CommandValues: commandValues,
// 			}
// 			driver.AsyncCh <- asyncValues

// 		}

// 	}

// }
