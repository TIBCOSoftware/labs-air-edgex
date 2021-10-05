package main

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"

	inttransforms "github.com/TIBCOSoftware/labs-air/edgexfoundry/app-service-inference-export-mqtt/internal/transforms"
	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

const (
	serviceKey = "app-service-inference-export-mqtt"
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
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	AccessToken string  `json:"accessToken"`
	Createdts   int64   `json:"createdts"`
	Updatedts   int64   `json:"updatedts"`
}

type metadataStruct struct {
	Gateway gatewayStruct   `json:"gateway"`
	Devices []models.Device `json:"devices"`
}

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
	appSettings := edgexSdk.ApplicationSettings()
	if appSettings != nil {

		gatewayID := inttransforms.GetAppSetting(appSettings, "GatewayId")
		gatewayDescription := inttransforms.GetAppSetting(appSettings, "GatewayDescription")
		gatewayHostname := inttransforms.GetAppSetting(appSettings, "GatewayHostname")
		gatewayLatitude := inttransforms.GetAppSetting(appSettings, "GatewayLatitude")
		gatewayLongitude := inttransforms.GetAppSetting(appSettings, "GatewayLongitude")
		gatewayAccessToken := inttransforms.GetAppSetting(appSettings, "GatewayAccessToken")

		startupTime := time.Now().UnixNano() / int64(time.Millisecond)
		gatewayInfo.UUUID = gatewayID
		gatewayInfo.AccessToken = gatewayAccessToken
		gatewayInfo.Address = gatewayHostname
		gatewayInfo.Createdts = startupTime
		gatewayInfo.Updatedts = startupTime
		gatewayInfo.Description = gatewayDescription
		gatewayInfo.Latitude, _ = strconv.ParseFloat(gatewayLatitude, 32)
		gatewayInfo.Longitude, _ = strconv.ParseFloat(gatewayLongitude, 32)
	} else {
		fmt.Println("Application settings nil")
	}

	// Add HTTP Route for ModelConf Dynamic Configuration
	edgexSdk.AddRoute("/api/v1/addModelConf", processAddModelConfig, "POST")
	edgexSdk.AddRoute("/api/v1/deleteModelConf", processDeleteModelConfig, "POST")

	// Create the MQTT Sender
	mqttConfig, _ := inttransforms.LoadMQTTConfig(appSettings)
	mqttSender = inttransforms.NewMQTTSender(edgexSdk.LoggingClient, nil, mqttConfig)

	inttransforms.InitializeModeling(appSettings)

	// Set pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	edgexSdk.SetFunctionsPipeline(
		// transforms.NewFilter(deviceNames).FilterByDeviceName,
		processEvent,
	)

	edgexSdk.LoggingClient.Info(fmt.Sprintf("Set the function pipeline: %s", "Ready to run"))

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

	edgexcontext.LoggingClient.Info(fmt.Sprintf("Processing Event: %s", "Starting"))
	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}

	edgexcontext.LoggingClient.Debug(fmt.Sprintf("Event: %s", params[0].(models.Event)))

	if event, ok := params[0].(models.Event); ok {

		edgexcontext.LoggingClient.Info(fmt.Sprintf("Check if inference available for device: %s", event.Device))

		// Process only 1 reading.  Need to modify to process multiple readings

		// Check to see if resource is inferable
		if inttransforms.IsInferable(event.Readings[0].Device, event.Readings[0].Name) {

			edgexcontext.LoggingClient.Info(fmt.Sprintf("Inferencing for device: %s", event.Device))

			prediction, _ := inttransforms.Predict(event.Readings[0])
			edgexcontext.LoggingClient.Info(fmt.Sprintf("Prediction: %s", prediction))

			// Update Reading
			event.Readings[0].Id = uuid.New().String()
			event.Readings[0].ValueType = "String"
			event.Readings[0].Value = prediction
			event.Readings[0].Name = event.Readings[0].Name + "_Inferred"
			event.Readings[0].BinaryValue = nil

			// Create reading to be exported
			jsondat := &msgStruct{
				ID:       event.ID,
				Device:   event.Device,
				Origin:   event.Origin,
				Gateway:  gatewayInfo.UUUID,
				Readings: event.Readings,
			}

			encjson, _ := json.Marshal(jsondat)

			edgexcontext.LoggingClient.Info(fmt.Sprintf("New Event: %s", encjson))

			// Export event
			return mqttSender.MQTTSend(string(encjson))

		}

		return false, nil

	}

	return false, errors.New("Unexpected type received")

}

func processAddModelConfig(writer http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		fmt.Printf("Processing addModelConfig request ERROR\n")
	}

	inttransforms.RegisterModel(body)

	writer.Header().Set("Content-Type", "text/plain")
	writer.Write([]byte("success config"))
}

func processDeleteModelConfig(writer http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		fmt.Printf("Processing deleteRule request ERROR\n")
	}

	inttransforms.UnregisterModel(body)

	writer.Header().Set("Content-Type", "text/plain")
	writer.Write([]byte("success config"))
}
