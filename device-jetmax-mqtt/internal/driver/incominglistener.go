// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	name = "name"
	cmd  = "cmd"
)

type BoolData struct {
	Data bool `json:"data"`
}

type IntData struct {
	Data int `json:"data"`
}


func (d *Driver) onIncomingDataReceived(client mqtt.Client, message mqtt.Message) {
	driver.Logger.Debug(fmt.Sprintf("[Incoming listener] Incoming reading received: topic=%v payload=%v", message.Topic(), string(message.Payload())))

	var deviceName string
	var resourceName string
	var reading interface{}


	incomingTopic := message.Topic()
	subscribedTopic := d.serviceConfig.MQTTBrokerInfo.IncomingTopic
	subscribedTopic = strings.Replace(subscribedTopic, "#", "", -1)
	incomingTopic = strings.Replace(incomingTopic, subscribedTopic, "", -1)

	deviceName = "JetmaxDevice"

	var data map[string]interface{}
	err := json.Unmarshal(message.Payload(), &data)
	if err != nil {
		driver.Logger.Errorf("Error unmarshaling payload: %s", err)
		return
	}

	driver.Logger.Debug(fmt.Sprintf("Incoming data: %+v", data))

	var ok bool
	switch incomingTopic {
	case "arm/status":
		resourceName = "servo1"
		reading, ok = data[resourceName]
		if !ok {
			driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No reading data found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
			return
		}
		publishReading(deviceName, resourceName, reading)

		resourceName = "servo2"
		reading, ok = data[resourceName]
		if !ok {
			driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No reading data found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
			return
		}
		publishReading(deviceName, resourceName, reading)

		resourceName = "servo3"
		reading, ok = data[resourceName]
		if !ok {
			driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No reading data found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
			return
		}
		publishReading(deviceName, resourceName, reading)

		resourceName = "joint1"
		reading, ok = data[resourceName]
		if !ok {
			driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No reading data found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
			return
		}
		publishReading(deviceName, resourceName, reading)

		resourceName = "joint2"
		reading, ok = data[resourceName]
		if !ok {
			driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No reading data found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
			return
		}
		publishReading(deviceName, resourceName, reading)

		resourceName = "joint3"
		reading, ok = data[resourceName]
		if !ok {
			driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No reading data found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
			return
		}
		publishReading(deviceName, resourceName, reading)
	default:
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No reading data found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
		return
	}


}


func publishReading(deviceName string, resourceName string, reading interface{}) {

	service := service.RunningService()

	deviceObject, ok := service.DeviceResource(deviceName, resourceName)
	if !ok {
		driver.Logger.Errorf("[Incoming listener] Incoming reading ignored, device resource `%s` not found from the device `%s`", resourceName, deviceName)
		return
	}

	req := models.CommandRequest{
		DeviceResourceName: resourceName,
		Type:               deviceObject.Properties.ValueType,
	}

	result, err := newResult(req, reading)
	if err != nil {
		driver.Logger.Errorf("[Incoming listener] Incoming reading ignored, %v", err)
		return
	}

	asyncValues := &models.AsyncValues{
		DeviceName:    deviceName,
		CommandValues: []*models.CommandValue{result},
	}

	driver.AsyncCh <- asyncValues

}