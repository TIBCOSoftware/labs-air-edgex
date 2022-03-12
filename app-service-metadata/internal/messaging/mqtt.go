package messaging

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/labs-air/edgexfoundry/app-service-metadata/internal/config"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/util"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
)

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
	opts   config.MQTTConfig
	lc     logger.LoggingClient
}

// NewMQTTSender creates, initializes and returns a new instance of MQTTSender
func NewMQTTSender(lc logger.LoggingClient, keyCertPair *KeyCertPair, mqttConfig *config.MQTTConfig) *MQTTSender {
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
			lc.Error("Failed loading x509 data")
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
		lc.Warn(fmt.Sprintf("Connection lost : %v", e))
		token := client.Connect()
		if token.Wait() && token.Error() != nil {
			lc.Info(fmt.Sprintf("Reconnection failed : %v", token.Error()))
		} else {
			lc.Info("Reconnection sucessful")
		}
	})

	client := MQTT.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		lc.Error("Failed to connect: %v", token.Error())
		return nil
	}

	lc.Info("Connected to MQTT sucessfully")

	sender := &MQTTSender{
		client: client,
		topic:  mqttConfig.Topic,
		opts:   *mqttConfig,
		lc:     lc,
	}

	return sender
}

func (sender MQTTSender) MQTTSend(msg string) (bool, interface{}) {

	sender.lc.Debug(fmt.Sprintf("Sending Message: [%s]\n", msg))

	if !sender.client.IsConnected() {
		sender.lc.Info("Reconnecting to mqtt server")
		if token := sender.client.Connect(); token.Wait() && token.Error() != nil {
			return false, fmt.Errorf("could not connect to mqtt server, drop event. Error: %s", token.Error().Error())
		}
		sender.lc.Info("Reconnected to mqtt server")
	}

	data, err := util.CoerceType(msg)
	if err != nil {
		return false, err
	}

	token := sender.client.Publish(sender.topic, sender.opts.QOS, sender.opts.Retain, data)
	// FIXME: could be removed? set of tokens?
	token.Wait()

	if token.Error() != nil {
		return false, token.Error()
	}

	// loggingClient.Trace("Data exported", "Transport", "MQTT", clients.CorrelationHeader)

	return true, nil
}

func (sender MQTTSender) SendMetadata(msg string) (bool, interface{}) {

	sender.lc.Trace(fmt.Sprintf("Sending Message: [%s]\n", msg))

	if !sender.client.IsConnected() {
		sender.lc.Info("Reconnecting to mqtt server")
		if token := sender.client.Connect(); token.Wait() && token.Error() != nil {
			return false, fmt.Errorf("could not connect to mqtt server, drop event. Error: %s", token.Error().Error())
		}
		sender.lc.Info("Reconnected to mqtt server")
	}

	data, err := util.CoerceType(msg)
	if err != nil {
		return false, err
	}

	token := sender.client.Publish("EdgexGatewayMetadata", sender.opts.QOS, sender.opts.Retain, data)
	// FIXME: could be removed? set of tokens?
	token.Wait()

	if token.Error() != nil {
		return false, token.Error()
	}

	// loggingClient.Trace("Data exported", "Transport", "MQTT", clients.CorrelationHeader)

	return true, nil
}
