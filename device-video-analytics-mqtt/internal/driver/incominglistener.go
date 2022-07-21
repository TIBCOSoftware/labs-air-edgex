// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
)


func (d *Driver) onIncomingDataReceived(client mqtt.Client, message mqtt.Message) {
	driver.Logger.Debug(fmt.Sprintf("[Incoming listener] Incoming reading received: topic=%v msg=%v", message.Topic(), string(message.Payload())))

	var deviceName string
	var resourceName string

	incomingTopic := message.Topic()
	subscribedTopic := d.serviceConfig.MQTTBrokerInfo.IncomingTopic
	subscribedTopic = strings.Replace(subscribedTopic, "#", "", -1)
	incomingTopic = strings.Replace(incomingTopic, subscribedTopic, "", -1)

	var data map[string]interface{}
	json.Unmarshal(message.Payload(), &data)

	driver.Logger.Debug(fmt.Sprintf("Incoming data: %+v", data))

	// if !checkDataWithKey(data, "deviceName") || !checkDataWithKey(data, "resourceName") {
	// 	return
	// }

	// deviceName := data["deviceName"].(string)
	// resourceName := data["resourceName"].(string)

	deviceName = "VADevice"

	switch incomingTopic {
	case "event/alert":
		resourceName = "model_score"
		reading := data

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
