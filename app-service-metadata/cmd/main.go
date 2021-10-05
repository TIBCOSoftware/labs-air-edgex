package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	inttransforms "github.com/TIBCOSoftware/labs-air/edgexfoundry/app-service-metadata/internal/transforms"
	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/google/uuid"
)

const (
	serviceKey = "app-service-metadata"
)

type msgStruct struct {
	ID       string           `json:"id"`
	Device   string           `json:"device"`
	Origin   int64            `json:"source"`
	Gateway  string           `json:"gateway"`
	Readings []models.Reading `json:"readings"`
}

type gatewayStruct struct {
	UUUID       string  `json:"uuid"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	Router      string  `json:"router"`
	RouterPort  string  `json:"routerPort"`
	DeployNetwork  string  `json:"deployNetwork"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	AccessToken string  `json:"accessToken"`
	Username    string  `json:"username"`
	Platform    string  `json:"platform"`
	Createdts   int64   `json:"createdts"`
	Updatedts   int64   `json:"updatedts"`
}

type metadataStruct struct {
	Gateway gatewayStruct   `json:"gateway"`
	Devices []models.Device `json:"devices"`
}

var mdc metadata.DeviceClient
var mqttSender *inttransforms.MQTTSender
var gatewayInfo gatewayStruct

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
	metadataPublishInterval := 30
	metadataClient := ""
	appSettings := edgexSdk.ApplicationSettings()
	if appSettings != nil {

		gatewayID := inttransforms.GetAppSetting(appSettings, "GatewayId")
		gatewayDescription := inttransforms.GetAppSetting(appSettings, "GatewayDescription")
		gatewayHostname := inttransforms.GetAppSetting(appSettings, "GatewayHostname")
		gatewayRouter := inttransforms.GetAppSetting(appSettings, "GatewayRouter")
		gatewayRouterPort := inttransforms.GetAppSetting(appSettings, "GatewayRouterPort")
		gatewayDeployNetwork := inttransforms.GetAppSetting(appSettings, "GatewayDeployNetwork")
		gatewayLatitude := inttransforms.GetAppSetting(appSettings, "GatewayLatitude")
		gatewayLongitude := inttransforms.GetAppSetting(appSettings, "GatewayLongitude")
		gatewayAccessToken := inttransforms.GetAppSetting(appSettings, "GatewayAccessToken")
		gatewayUsername := inttransforms.GetAppSetting(appSettings, "GatewayUsername")
		gatewayPlatform := inttransforms.GetAppSetting(appSettings, "GatewayPlatform")
		metadataClient = inttransforms.GetAppSetting(appSettings, "MetadataClient")
		metadataPublishIntervalStr := inttransforms.GetAppSetting(appSettings, "MetadataPublishIntervalSecs")

		startupTime := time.Now().UnixNano() / int64(time.Millisecond)
		gatewayInfo.UUUID = gatewayID
		gatewayInfo.AccessToken = gatewayAccessToken
		gatewayInfo.Address = gatewayHostname
		gatewayInfo.Router = gatewayRouter
		gatewayInfo.RouterPort = gatewayRouterPort
		gatewayInfo.DeployNetwork = gatewayDeployNetwork
		gatewayInfo.Createdts = startupTime
		gatewayInfo.Updatedts = startupTime
		gatewayInfo.Description = gatewayDescription
		gatewayInfo.Username = gatewayUsername
		gatewayInfo.Platform = gatewayPlatform
		gatewayInfo.Latitude, _ = strconv.ParseFloat(gatewayLatitude, 32)
		gatewayInfo.Longitude, _ = strconv.ParseFloat(gatewayLongitude, 32)
		metadataPublishInterval64, _ := strconv.ParseInt(metadataPublishIntervalStr, 10, 32)
		metadataPublishInterval = int(metadataPublishInterval64)

	} else {
		fmt.Println("Application settings nil")
	}

	deviceClientURL := metadataClient + clients.ApiDeviceRoute

	mdc = metadata.NewDeviceClient(
		local.New(deviceClientURL))
	// local.New(config.Clients[common.CoreDataClientName].Url() + clients.ApiValueDescriptorRoute))

	// Create the MQTT Sender
	mqttConfig, _ := inttransforms.LoadMQTTConfig(appSettings)
	mqttSender = inttransforms.NewMQTTSender(edgexSdk.LoggingClient, nil, mqttConfig)

	// Set pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	edgexSdk.SetFunctionsPipeline(
		// transforms.NewFilter(deviceNames).FilterByDeviceName,
		processEvent,
	)

	// Set ticker to trigger metadata publishing
	ticker := time.NewTicker(time.Duration(metadataPublishInterval) * time.Second)
	go func() {
		for ; true; <-ticker.C {
			publishMetadata(edgexSdk.LoggingClient)
		}
	}()

	// Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err := edgexSdk.MakeItRun()
	if err != nil {
		edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	ticker.Stop()
	os.Exit(0)
}

func processEvent(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {

	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}

	if event, ok := params[0].(models.Event); ok {

		edgexcontext.LoggingClient.Debug(fmt.Sprintf("Processing event for device: %s", event.Device))

		// Check to see if reading is of type binary and if so, update reading
		if event.Readings[0].ValueType == "Binary" {
			edgexcontext.LoggingClient.Debug(fmt.Sprintf("Processing Binary Data"))
			event.Readings[0].Id = uuid.New().String()
			// event.Readings[0].Value = string(event.Readings[0].BinaryValue)
			event.Readings[0].Value = base64.StdEncoding.EncodeToString(event.Readings[0].BinaryValue)
			// event.Readings[0].Value = "image"
			event.Readings[0].BinaryValue = nil
			event.Readings[0].ValueType = "String"
		}

		jsondat := &msgStruct{
			ID:       event.ID,
			Device:   event.Device,
			Origin:   event.Origin,
			Gateway:  gatewayInfo.UUUID,
			Readings: event.Readings,
		}

		encjson, _ := json.Marshal(jsondat)

		edgexcontext.Complete(encjson)

		// return false, nil

		return true, nil
	}

	return false, errors.New("Unexpected type received")

}

func publishMetadata(loggingClient logger.LoggingClient) (bool, interface{}) {

	loggingClient.Debug("Publishing metadata")

	// Update gatewayInfo with publish time
	publishTime := time.Now().UnixNano() / int64(time.Millisecond)
	gatewayInfo.Updatedts = publishTime

	devices, err := mdc.Devices(context.Background())

	if devices != nil && err == nil {
		jsonmd := &metadataStruct{
			Gateway: gatewayInfo,
			Devices: devices,
		}

		encjson, _ := json.Marshal(jsonmd)

		// Export metadata
		return mqttSender.MQTTSendMetadata(string(encjson))
	}

	return false, errors.New("Unexpected devices result received")
}
