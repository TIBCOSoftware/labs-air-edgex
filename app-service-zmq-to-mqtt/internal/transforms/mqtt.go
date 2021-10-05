package transforms

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	// "github.com/edgexfoundry/app-functions-sdk-go/pkg/util"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

// MQTTConfig contains mqtt client parameters
type MQTTConfig struct {
	Protocol      string
	Hostname      string
	Port          string
	TrustStore    string
	User          string
	Password      string
	Publisher     string
	Topic         string
	QOS           byte
	Retain        bool
	AutoReconnect bool
}

// KeyCertPair is used to pass key/cert pair to NewMQTTSender
// KeyPEMBlock and CertPEMBlock will be used if they are not nil
// then it will fall back to KeyFile and CertFile
type KeyCertPair struct {
	KeyFile      string
	CertFile     string
	KeyPEMBlock  []byte
	CertPEMBlock []byte
}

// MQTTSender ...
type MQTTSender struct {
	client MQTT.Client
	topic  string
	opts   MQTTConfig
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

// LoadMQTTConfig Loads the MQTT configuration necessary to connect to MQTT
func LoadMQTTConfig(appSettings map[string]string) (*MQTTConfig, error) {

	var protocol, host, port, publisher, user, password, trustStore, topic string

	if appSettings != nil {
		protocol = GetAppSetting(appSettings, "MqttProtocol")
		host = GetAppSetting(appSettings, "MqttHostname")
		port = GetAppSetting(appSettings, "MqttPort")
		publisher = GetAppSetting(appSettings, "MqttPublisher")
		user = GetAppSetting(appSettings, "MqttUser")
		password = GetAppSetting(appSettings, "MqttPassword")
		trustStore = GetAppSetting(appSettings, "MqttTrustStore")
		topic = GetAppSetting(appSettings, "MqttTopic")
	} else {
		return nil, errors.New("No application-specific settings found")
	}

	config := MQTTConfig{}

	config.Protocol = protocol
	config.Hostname = host
	config.Port = port
	config.Publisher = publisher
	config.User = user
	config.Password = password
	config.TrustStore = trustStore
	config.Topic = topic
	config.QOS = 0
	config.Retain = false
	config.AutoReconnect = false

	return &config, nil
}

// NewMQTTSender creates, initializes and returns a new instance of MQTTSender
func NewMQTTSender(logging logger.LoggingClient, keyCertPair *KeyCertPair, mqttConfig *MQTTConfig) *MQTTSender {
	protocol := strings.ToLower(mqttConfig.Protocol)

	opts := MQTT.NewClientOptions()
	broker := protocol + "://" + mqttConfig.Hostname + ":" + mqttConfig.Port
	opts.AddBroker(broker)
	opts.SetClientID(mqttConfig.Publisher)
	opts.SetUsername(mqttConfig.User)
	opts.SetPassword(mqttConfig.Password)
	// opts.SetAutoReconnect(mqttConfig.AutoReconnect)
	opts.SetAutoReconnect(true)
	// opts.SetKeepAlive(time.Second * time.Duration(30))

	if (protocol == "tcps" || protocol == "ssl" || protocol == "tls") && keyCertPair != nil {
		var cert tls.Certificate
		var err error

		if keyCertPair.KeyPEMBlock != nil && keyCertPair.CertPEMBlock != nil {
			cert, err = tls.X509KeyPair(keyCertPair.CertPEMBlock, keyCertPair.KeyPEMBlock)
		} else {
			cert, err = tls.LoadX509KeyPair(keyCertPair.CertFile, keyCertPair.KeyFile)
		}

		if err != nil {
			logging.Error("Failed loading x509 data")
			return nil
		}

		tlsConfig := &tls.Config{
			ClientCAs:          nil,
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert},
		}

		opts.SetTLSConfig(tlsConfig)

	}

	opts.SetConnectionLostHandler(func(client MQTT.Client, e error) {
		logging.Warn(fmt.Sprintf("Connection lost : %v", e))
		token := client.Connect()
		if token.Wait() && token.Error() != nil {
			logging.Info(fmt.Sprintf("Reconnection failed : %v", token.Error()))
		} else {
			logging.Info(fmt.Sprintf("Reconnection sucessful"))
		}
	})

	client := MQTT.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		logging.Error("Failed to connect: %v", token.Error())
		return nil
	}

	logging.Info(fmt.Sprintf("Connected to MQTT sucessfully"))

	sender := &MQTTSender{
		client: client,
		topic:  mqttConfig.Topic,
		opts:   *mqttConfig,
	}

	// sender := &MQTTSender{
	// 	client: MQTT.NewClient(opts),
	// 	topic:  mqttConfig.Topic,
	// 	opts:   *mqttConfig,
	// }

	return sender
}

// MQTTSend will send data from the previous function to the specified Endpoint via MQTT.
// If no previous function exists, then the event that triggered the pipeline will be used.
// An empty string for the mimetype will default to application/json.
func (sender MQTTSender) MQTTSend(msg []byte) (bool, interface{}) {

	loggingClient.Debug(fmt.Sprintf("Sending Message: %s\n", msg))

	if !sender.client.IsConnected() {
		loggingClient.Info("Reconnecting to mqtt server")
		if token := sender.client.Connect(); token.Wait() && token.Error() != nil {
			return false, fmt.Errorf("Could not connect to mqtt server, drop event. Error: %s", token.Error().Error())
		}
		loggingClient.Info("Reconnected to mqtt server")
	}

	// data, err := util.CoerceType(msg)
	// if err != nil {
	// 	return false, err
	// }

	token := sender.client.Publish(sender.topic, sender.opts.QOS, sender.opts.Retain, msg)
	// FIXME: could be removed? set of tokens?
	token.Wait()

	if token.Error() != nil {
		return false, token.Error()
	}

	loggingClient.Trace("Data exported", "Transport", "MQTT", clients.CorrelationHeader)

	return true, nil
}

