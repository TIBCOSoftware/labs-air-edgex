// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
)

type BoolData struct {
	Data bool `json:"data"`
}

type IntData struct {
	Data int `json:"data"`
}


func onIncomingDataReceived(client mqtt.Client, message mqtt.Message) {
	driver.Logger.Debug(fmt.Sprintf("[Incoming listener] Incoming reading received: topic=%v payload=%v", message.Topic(), string(message.Payload())))

	var data map[string]interface{}
	json.Unmarshal(message.Payload(), &data)

	driver.Logger.Debug(fmt.Sprintf("Incoming data: %+v", data))

	driver.Logger.Debug(fmt.Sprintf("Incoming data value: %+v", data["data"]))

	if !checkDataWithKey(data, "deviceName") || !checkDataWithKey(data, "resourceName") {
		return
	}

	var resourceName string
	deviceName := "JetmaxDevice"
	// topic := message.Topic()

	switch message.Topic() {
	case "jetmax/servo1/status":
		resourceName = "servo1"
	default:
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No reading data found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
		return
	}

	// deviceName := data["deviceName"].(string)
	// resourceName := data["resourceName"].(string)

	reading, ok := data[resourceName]
	if !ok {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No reading data found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
		return
	}

	service := service.RunningService()

	deviceObject, ok := service.DeviceResource(deviceName, resourceName)
	if !ok {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No DeviceObject found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
		return
	}

	req := models.CommandRequest{
		DeviceResourceName: resourceName,
		Type:               deviceObject.Properties.ValueType,
	}

	result, err := newResult(req, reading)

	if err != nil {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored.   topic=%v msg=%v error=%v", message.Topic(), string(message.Payload()), err))
		return
	}

	asyncValues := &models.AsyncValues{
		DeviceName:    deviceName,
		CommandValues: []*models.CommandValue{result},
	}

	driver.Logger.Info(fmt.Sprintf("[Incoming listener] Incoming reading received: topic=%v msg=%v", message.Topic(), string(message.Payload())))

	driver.AsyncCh <- asyncValues

}

func checkDataWithKey(data map[string]interface{}, key string) bool {
	val, ok := data[key]
	if !ok {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No %v found : msg=%v", key, data))
		return false
	}

	switch val.(type) {
	case string:
		return true
	default:
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. %v should be string : msg=%v", key, data))
		return false
	}
}
