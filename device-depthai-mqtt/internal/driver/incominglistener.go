// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	sdkModel "github.com/edgexfoundry/device-sdk-go/pkg/models"
	sdk "github.com/edgexfoundry/device-sdk-go/pkg/service"
)

type ReportData struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
	Proximity   int     `json:"proximity"`
	Accx        int     `json:"acc_x"`
	Accy        int     `json:"acc_y"`
	Accz        int     `json:"acc_z"`
	Gyrx        int     `json:"gyr_x"`
	Gyry        int     `json:"gyr_y"`
	Gyrz        int     `json:"gyr_z"`
	Magx        int     `json:"mag_x"`
	Magy        int     `json:"mag_y"`
	Magz        int     `json:"mag_z"`
	Ts          int64   `json:"ts"`
	Mac         string  `json:"mac"`
}

type Report struct {
	Reported ReportData `json:"reported"`
}

type Stmsg struct {
	State Report `json:"state"`
}

type DepthAIMsg struct {
	Predict string `json:"predict"`
	HasImage bool  `json:"hasimage"`
	Image   []byte `json:"image"`
	Ts      int64  `json:"ts"`
}

func startIncomingListening() error {
	var scheme = driver.Config.IncomingSchema
	var brokerUrl = driver.Config.IncomingHost
	var brokerPort = driver.Config.IncomingPort
	var username = driver.Config.IncomingUser
	var password = driver.Config.IncomingPassword
	var mqttClientId = driver.Config.IncomingClientId
	var qos = byte(driver.Config.IncomingQos)
	var keepAlive = driver.Config.IncomingKeepAlive
	var topic = driver.Config.IncomingTopic

	uri := &url.URL{
		Scheme: strings.ToLower(scheme),
		Host:   fmt.Sprintf("%s:%d", brokerUrl, brokerPort),
		User:   url.UserPassword(username, password),
	}

	driver.Logger.Info(fmt.Sprintf("startIncomingListening calling createClient for qos: %v  and topic: %s", qos, topic))

	// client, err := createClient(mqttClientId, uri, keepAlive)
	// if err != nil {
	// 	return err
	// }

	var client mqtt.Client
	var err error
	for i := 1; i <= 10; i++ {
		client, err = createClient(mqttClientId, uri, keepAlive)
		if err != nil && i == 10 {
			return err
		} else if err != nil {
			driver.Logger.Error(fmt.Sprintf("Fail to initial conn for incoming data, %v ", err))
			time.Sleep(time.Duration(5) * time.Second)
			driver.Logger.Warn("Retry to initial conn for incoming data")
			continue
		}
		driver.Logger.Info("Created client successfully")
		break
	}

	defer func() {
		if client.IsConnected() {
			client.Disconnect(5000)
		}
	}()

	driver.Logger.Info(fmt.Sprintf("Subscribing for qos: %v  and topic: %s", qos, topic))

	token := client.Subscribe(topic, qos, onIncomingDataReceived)
	if token.Wait() && token.Error() != nil {
		driver.Logger.Info(fmt.Sprintf("[Incoming listener] Stop incoming data listening. Cause:%v", token.Error()))
		return token.Error()
	}

	driver.Logger.Info("[Incoming listener] Start incoming data listening. ")
	select {}
}

func onIncomingDataReceived(client mqtt.Client, message mqtt.Message) {
	driver.Logger.Debug(fmt.Sprintf("[Incoming listener] Incoming reading received: topic=%v msg=%v", message.Topic(), string(message.Payload())))

	data := DepthAIMsg{}
	json.Unmarshal(message.Payload(), &data)

	driver.Logger.Debug(fmt.Sprintf("Incoming data: %+v", data))

	// deviceName := "ST_" + data.State.Reported.Mac
	deviceName := "DepthAI_" + "OAK-D"
	deviceResource := ""
	tms := data.Ts

	service := sdk.RunningService()

	// Process predict reading
	deviceResource = "DAI_Snapshot_Inferred"
	deviceObject, ok := service.DeviceResource(deviceName, deviceResource, "get")
	if !ok {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No DeviceObject found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
		return
	}

	req := sdkModel.CommandRequest{
		DeviceResourceName: deviceResource,
		Type:               sdkModel.ParseValueType(deviceObject.Properties.Value.Type),
	}

	result, err := newResult(req, data.Predict, tms)

	if err != nil {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored.   topic=%v msg=%v error=%v", message.Topic(), string(message.Payload()), err))
		return
	}

	asyncValues := &sdkModel.AsyncValues{
		DeviceName:    deviceName,
		CommandValues: []*sdkModel.CommandValue{result},
	}

	driver.AsyncCh <- asyncValues

	// Process snapshot reading
	deviceResource = "DAI_Snapshot"
	deviceObject, ok = service.DeviceResource(deviceName, deviceResource, "get")
	if !ok {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No DeviceObject found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
		return
	}

	// Skip sending image if value is NONE
	if !data.HasImage {
		return
	}

	req = sdkModel.CommandRequest{
		DeviceResourceName: deviceResource,
		Type:               sdkModel.ParseValueType(deviceObject.Properties.Value.Type),
	}

	result, err = newResult(req, data.Image, tms)

	if err != nil {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored.   topic=%v msg=%v error=%v", message.Topic(), string(message.Payload()), err))
		return
	}

	asyncValues = &sdkModel.AsyncValues{
		DeviceName:    deviceName,
		CommandValues: []*sdkModel.CommandValue{result},
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
