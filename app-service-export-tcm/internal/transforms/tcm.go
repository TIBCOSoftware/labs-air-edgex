package transforms

import (
	"errors"
	"fmt"
	"log"

	"github.com/TIBCOSoftware/labs-air/edgexfoundry/app-service-export-tcm/internal/tibco.com/eftl"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

type TCMConfig struct {
	Url               string
	AuthenticationKey string
	ClientId          string
	Topic             string
}

// TCMSender ...
type TCMSender struct {
	conn  *eftl.Connection
	topic string
}

var loggingClient logger.LoggingClient

// SetLoggingClient set the logging client
func SetLoggingClient(logging logger.LoggingClient) {
	loggingClient = logging
}

// GetAppSetting get application setting
func GetAppSetting(settings map[string]string, name string) string {
	value, ok := settings[name]

	if ok {
		loggingClient.Info(fmt.Sprintf("Setting for %s: %s", name, value))
		return value
	}
	loggingClient.Error(fmt.Sprintf("ApplicationName application setting %s not found", name))
	return ""

}

// LoadTCMConfig Loads the tcm configuration necessary to connect to TCM
func LoadTCMConfig(sdk *appsdk.AppFunctionsSDK) (*TCMConfig, error) {
	if sdk == nil {
		return nil, errors.New("Invalid AppFunctionsSDK")
	}

	loggingClient = sdk.LoggingClient

	var url, authenticationKey, clientId, topic string

	appSettings := sdk.ApplicationSettings()
	if appSettings != nil {
		url = GetAppSetting(appSettings, "TCMUrl")
		authenticationKey = GetAppSetting(appSettings, "AuthenticationKey")
		clientId = GetAppSetting(appSettings, "ClientId")
		topic = GetAppSetting(appSettings, "TCMTopic")
	} else {
		return nil, errors.New("No application-specific settings found")
	}

	config := TCMConfig{}

	config.Url = url
	config.AuthenticationKey = authenticationKey
	config.ClientId = clientId
	config.Topic = topic

	return &config, nil
}

// NewTCMSender creates, initializes and returns a new instance of TCMSender
func NewTCMSender(logging logger.LoggingClient, config *TCMConfig) *TCMSender {

	// Channel on which to receive connection errors.
	errChan := make(chan error, 1)

	// Set connection options.
	opts := &eftl.Options{
		Password: config.AuthenticationKey, ClientID: config.ClientId,
	}

	// Connect to TIBCO Cloud Messaging.
	conn, err := eftl.Connect(config.Url, opts, errChan)
	if err != nil {
		loggingClient.Error(fmt.Sprintf("EFTL Connect got error %s", err.Error()))
		return nil
	}

	// Listen for connection errors.
	go func() {
		for err := range errChan {
			loggingClient.Error(fmt.Sprintf("EFTL connection error %s", err.Error()))
			log.Println("connection error:", err)
		}
	}()

	return &TCMSender{
		conn:  conn,
		topic: config.Topic,
	}
}

// TCMSend will send data from the previous function to the specified Endpoint via TCM.
// If no previous function exists, then the event that triggered the pipeline will be used.
// An empty string for the mimetype will default to application/json.
func (sender TCMSender) TCMSend(msg string) (bool, interface{}) {

	loggingClient.Info(fmt.Sprintf("TCMSend event: "))

	eftlMsg := eftl.Message{
		"topic":   sender.topic,
		"jsonstr": msg,
	}

	// Publish the message to TIBCO Cloud Messaging.
	err := sender.conn.Publish(eftlMsg)

	if err != nil {
		loggingClient.Info(fmt.Sprintf("Error sending event: "))
		return false, fmt.Errorf("failed to send TCM message for reason [%s]", err.Error())
	}

	loggingClient.Debug(fmt.Sprintf("TCM message [%v] sent successfully", msg))

	return true, nil

}

// TCMSend will send data from the previous function to the specified Endpoint via TCM.
// If no previous function exists, then the event that triggered the pipeline will be used.
// An empty string for the mimetype will default to application/json.
func (sender TCMSender) TCMSendMetadata(msg string) (bool, interface{}) {

	loggingClient.Info(fmt.Sprintf("TCMSend event: "))

	eftlMsg := eftl.Message{
		"topic":   "EdgexGatewayMetadata",
		"jsonstr": msg,
	}

	// Publish the message to TIBCO Cloud Messaging.
	err := sender.conn.Publish(eftlMsg)

	if err != nil {
		loggingClient.Info(fmt.Sprintf("Error sending event: "))
		return false, fmt.Errorf("failed to send TCM message for reason [%s]", err.Error())
	}

	loggingClient.Debug(fmt.Sprintf("TCM message [%v] sent successfully", msg))

	return true, nil

}
