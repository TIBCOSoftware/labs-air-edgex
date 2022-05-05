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

const (
	name = "tname"
)


func onIncomingDataReceived(client mqtt.Client, message mqtt.Message) {
	driver.Logger.Debug(fmt.Sprintf("[Incoming listener] Incoming reading received: topic=%v payload=%v", message.Topic(), string(message.Payload())))

	var deviceName string
	var resourceName string
	var reading interface{}

	var data map[string]interface{}
	err := json.Unmarshal(message.Payload(), &data)
	if err != nil {
		driver.Logger.Errorf("Error unmarshaling payload: %s", err)
		return
	}

	driver.Logger.Debug(fmt.Sprintf("Incoming data: %+v", data))


	nameVal, ok := data[name]
	if !ok {
		driver.Logger.Errorf("[Incoming listener] Incoming reading ignored, reading data `%v` should contain the field `%s` to indicate the device name", data, name)
		return
	}

	deviceName = fmt.Sprintf("%s", nameVal)
	resourceName = "vibration"

	reading = data

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
