package functions

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	clientint "github.com/edgexfoundry/go-mod-core-contracts/v2/clients/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/google/uuid"

	"github.com/TIBCOSoftware/labs-air/edgexfoundry/app-service-metadata/internal/config"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/app-service-metadata/internal/messaging"
)

var gatewayConfig *config.GatewayConfig

type deviceInfo struct {
	Id             string             `json:"id"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	AdminState     string             `json:"adminState"`
	OperatingState string             `json:"operatingState"`
	LastConnected  int64              `json:"lastConnected"`
	LastReported   int64              `json:"lastReported"`
	Labels         []string           `json:"labels"`
	Location       interface{}        `json:"location"`
	ServiceName    string             `json:"serviceName"`
	ProfileName    string             `json:"profileName"`
	Profile        dtos.DeviceProfile `json:"profile"`
}

type msgStruct struct {
	ID       string             `json:"id"`
	DeviceName   string             `json:"device"`
	Origin   int64              `json:"source"`
	Gateway  string             `json:"gateway"`
	Readings []dtos.BaseReading `json:"readings"`
}

type metadataStruct struct {
	Gateway *config.GatewayConfig `json:"gateway"`
	Devices []deviceInfo          `json:"devices"`
}

type Export struct {
	lc                  logger.LoggingClient
	appConfig           *config.AppConfig
	deviceClient        clientint.DeviceClient
	deviceProfileClient clientint.DeviceProfileClient
	sender              *messaging.MQTTSender
}

func NewExport(lc logger.LoggingClient, appConfig *config.AppConfig, deviceClient clientint.DeviceClient, deviceProfileClient clientint.DeviceProfileClient) Export {

	gatewayConfig, _ = appConfig.LoadGatewayConfig()
	mqttConfig, _ := appConfig.LoadMQTTConfig()
	sender := messaging.NewMQTTSender(lc, nil, mqttConfig)

	return Export{
		lc:                  lc,
		appConfig:           appConfig,
		deviceClient:        deviceClient,
		deviceProfileClient: deviceProfileClient,
		sender:              sender,
	}
}

func (exp Export) ProcessEvent(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {

	exp.lc.Debug("Processing event")

	if data == nil {
		// We didn't receive a result
		return false, errors.New("processEvent: No data received")
	}

	if event, ok := data.(dtos.Event); ok {

		exp.lc.Debug(fmt.Sprintf("Processing event for device: %s", event.DeviceName))

		// Check to see if reading is of type binary and if so, update reading
		if event.Readings[0].ValueType == "Binary" {
			exp.lc.Debug("Processing Binary Data")
			event.Readings[0].Id = uuid.New().String()
			// event.Readings[0].Value = string(event.Readings[0].BinaryValue)
			event.Readings[0].Value = base64.StdEncoding.EncodeToString(event.Readings[0].BinaryValue)
			// event.Readings[0].Value = "image"
			event.Readings[0].BinaryValue = nil
			event.Readings[0].ValueType = "String"
		}

		jsondat := &msgStruct{
			ID:			event.Id,
			DeviceName: event.DeviceName,
			Origin:   	event.Origin,
			Gateway:  	gatewayConfig.UUID,
			Readings: 	event.Readings,
		}

		encjson, _ := json.Marshal(jsondat)

		exp.lc.Debug(fmt.Sprintf("Publish event json: %s", string(encjson)))

		ctx.SetResponseData(encjson)

		return true, nil
	}

	return false, errors.New("unexpected type received")
}

func (exp Export) PublishMetadata() (bool, interface{}) {

	exp.lc.Debug("Publishing metadata")

	metadata, err := exp.getMetadata()

	if err == nil {

		encjson, _ := json.Marshal(metadata)

		exp.lc.Debug(fmt.Sprintf("Publish metadata devices json: %s", string(encjson)))

		// Export metadata
		return exp.sender.SendMetadata(string(encjson))
	}

	return false, err
}

func (exp Export) getDevicesInfo() ([]deviceInfo, error) {

	// Get all devices
	devicesRes, err := exp.deviceClient.AllDevices(context.Background(), nil, 0, -1)

	if err == nil {

		dcDevices := devicesRes.Devices

		numDevices := len(dcDevices)

		devices := make([]deviceInfo, numDevices)

		// For each device
		for i, dev := range dcDevices {

			// Get profile
			profileRes, _ := exp.deviceProfileClient.DeviceProfileByName(context.Background(), dev.ProfileName)

			devices[i].Id = dev.Id
			devices[i].Name = dev.Name
			devices[i].Description = dev.Description
			devices[i].AdminState = dev.AdminState
			devices[i].OperatingState = dev.OperatingState
			devices[i].LastConnected = dev.LastConnected
			devices[i].LastReported = dev.LastReported
			devices[i].Labels = dev.Labels
			devices[i].Location = dev.Location
			devices[i].ServiceName = dev.ServiceName
			devices[i].ProfileName = dev.ProfileName
			devices[i].Profile = profileRes.Profile

		}

		return devices, nil

	} else {
		return nil, err
	}

}

func (exp Export) getMetadata() (*metadataStruct, error) {

	devices, err := exp.getDevicesInfo()

	if err == nil {
		// Update gatewayConfig with publish time
		publishTime := time.Now().UnixNano() / int64(time.Millisecond)
		gatewayConfig.Updatedts = publishTime

		metadata := &metadataStruct{
			Gateway: gatewayConfig,
			Devices: devices,
		}
		return metadata, err
	}

	return nil, err

}
