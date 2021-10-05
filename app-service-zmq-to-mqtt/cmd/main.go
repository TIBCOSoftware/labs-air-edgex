package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"

	inttransforms "github.com/TIBCOSoftware/labs-air/edgexfoundry/app-service-zmq-to-mqtt/internal/transforms"
	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/google/uuid"
)

const (
	serviceKey = "app-service-zmq-to-mqtt"
)

type msgStruct struct {
	ID       string           `json:"id"`
	Device   string           `json:"device"`
	Origin   int64            `json:"source"`
	Gateway  string           `json:"gateway"`
	Readings []models.Reading `json:"readings"`
}

var mqttSender *inttransforms.MQTTSender
var gatewayID string

func main() {

	// Create an instance of the EdgeX SDK and initialize it.
	edgexSdk := &appsdk.AppFunctionsSDK{ServiceKey: serviceKey}
	if err := edgexSdk.Initialize(); err != nil {
		message := fmt.Sprintf("SDK initialization failed: %v\n", err)
		if edgexSdk.LoggingClient != nil {
			edgexSdk.LoggingClient.Error(message)
		} else {
			fmt.Println(message)
		}
		os.Exit(-1)
	}

	// Set the logging client for the transforms
	inttransforms.SetLoggingClient(edgexSdk.LoggingClient)

	// Get the application's specific configuration settings.
	appSettings := edgexSdk.ApplicationSettings()

	if appSettings != nil {

		gatewayID = inttransforms.GetAppSetting(appSettings, "GatewayId")
	}

	// Create the MQTT Sender
	mqttConfig, _ := inttransforms.LoadMQTTConfig(appSettings)
	mqttSender = inttransforms.NewMQTTSender(edgexSdk.LoggingClient, nil, mqttConfig)

	// Set pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	edgexSdk.SetFunctionsPipeline(
		// transforms.NewFilter(deviceNames).FilterByDeviceName,
		processEvent,
	)

	// Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err := edgexSdk.MakeItRun()
	if err != nil {
		edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}

func processEvent(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {

	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}

	edgexcontext.LoggingClient.Debug(fmt.Sprintf("Event: %s", params[0].(models.Event)))

	if event, ok := params[0].(models.Event); ok {

		extraValue := ""

		edgexcontext.LoggingClient.Debug(fmt.Sprintf("Processing event for device: %s", event.Device))

		// Check to see if reading is of type binary and if so, update reading
		if event.Readings[0].ValueType == "Binary" {
			event.Readings[0].Id = uuid.New().String()
			// event.Readings[0].Value = string(event.Readings[0].BinaryValue)
			event.Readings[0].Value = base64.StdEncoding.EncodeToString(event.Readings[0].BinaryValue)
			// event.Readings[0].Value = "image"
			event.Readings[0].BinaryValue = nil
			event.Readings[0].ValueType = "String"
		}

		for i, reading := range event.Readings {

			var value bytes.Buffer
			values := make(map[string]interface{})
			error := json.Unmarshal([]byte(reading.Value), &values)

			if nil == error {
				for aKey, aValue := range values {
					edgexcontext.LoggingClient.Info(fmt.Sprintf("Pair Name: %s Key: %s", reading.Name, aKey))
					if "cv-roi-event" == reading.Name && "frame" == aKey {
						edgexcontext.LoggingClient.Info(fmt.Sprintf("Skipping: %s", reading.Name))
						sValue := aValue.(string)
						extraValue = sValue[2 : len(sValue)-2]
						continue
					}

					value.WriteString(aKey)
					value.WriteString(":")
					value.WriteString(fmt.Sprintf("%v", aValue))
					value.WriteString(",")
				}
				event.Readings[i].Value = value.String()
			}

			event.Readings[i].Created = adjustTimestampUnit(event.Readings[i].Created)
			event.Readings[i].Origin = adjustTimestampUnit(event.Readings[i].Origin)
			event.Readings[i].Modified = adjustTimestampUnit(event.Readings[i].Modified)

		}

		jsondat := &msgStruct{
			ID:       event.ID,
			Device:   event.Device,
			Origin:   adjustTimestampUnit(event.Origin),
			Gateway:  gatewayID,
			Readings: event.Readings,
		}

		encjson, _ := json.Marshal(jsondat)

		edgexcontext.LoggingClient.Info(fmt.Sprintf("Cleaned Event: %s", string(encjson)))

		if extraValue != "" {
			// Export event
			mqttSender.MQTTSend(encjson)

			event.Readings[0].Id = uuid.New().String()
			event.Readings[0].Value = extraValue
			event.Readings[0].Name = event.Readings[0].Name + "_Image"

			encjson, _ := json.Marshal(jsondat)

			edgexcontext.LoggingClient.Info(fmt.Sprintf("Frame Event: %s", string(encjson)))
			return mqttSender.MQTTSend(encjson)

		} else {
			return mqttSender.MQTTSend(encjson)
		}

		// return false, nil
	}

	return false, errors.New("Unexpected type received")

}

func adjustTimestampUnit(timestamp int64) int64 {
	exp := math.Log10(float64(timestamp))
	return int64(float64(timestamp) * math.Pow10(12-int(exp)))
}
