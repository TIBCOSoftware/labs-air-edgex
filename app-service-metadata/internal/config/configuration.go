package config

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
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

type GatewayConfig struct {
	UUID          string  `json:"uuid"`
	Description   string  `json:"description"`
	Address       string  `json:"address"`
	Router        string  `json:"router"`
	RouterPort    string  `json:"routerPort"`
	DeployNetwork string  `json:"deployNetwork"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	AccessToken   string  `json:"accessToken"`
	Username      string  `json:"username"`
	Platform      string  `json:"platform"`
	Createdts     int64   `json:"createdts"`
	Updatedts     int64   `json:"updatedts"`
}

type AppConfig struct {
	lc          logger.LoggingClient
	appSettings map[string]string
}

func NewAppConfig(lc logger.LoggingClient, settings map[string]string) *AppConfig {

	if settings != nil {
		return &AppConfig{
			lc:          lc,
			appSettings: settings,
		}
	} else {
		return nil
	}
}

// GetAppSetting get application setting
func (ac AppConfig) GetAppSetting(name string) string {
	value, ok := ac.appSettings[name]

	if ok {
		ac.lc.Info(fmt.Sprintf("Setting for %s: %s", name, value))
		return value
	}
	ac.lc.Error(fmt.Sprintf("ApplicationName application setting %s not found", name))
	return ""

}

// LoadMQTTConfig Loads the MQTT configuration necessary to connect to MQTT
func (ac AppConfig) LoadMQTTConfig() (*MQTTConfig, error) {

	config := MQTTConfig{}

	if ac.appSettings != nil {
		config.Protocol = ac.GetAppSetting("MqttProtocol")
		config.Hostname = ac.GetAppSetting("MqttHostname")
		config.Port = ac.GetAppSetting("MqttPort")
		config.Publisher = ac.GetAppSetting("MqttPublisher")
		config.User = ac.GetAppSetting("MqttUser")
		config.Password = ac.GetAppSetting("MqttPassword")
		config.TrustStore = ac.GetAppSetting("MqttTrustStore")
		config.Topic = ac.GetAppSetting("MqttTopic")
		config.QOS = 0
		config.Retain = false
		config.AutoReconnect = false
	} else {
		return nil, errors.New("no application-specific settings found")
	}

	return &config, nil
}

// LoadMQTTConfig Loads the MQTT configuration necessary to connect to MQTT
func (ac AppConfig) LoadGatewayConfig() (*GatewayConfig, error) {

	config := GatewayConfig{}

	if ac.appSettings != nil {
		gatewayLatitude := ac.GetAppSetting("GatewayLatitude")
		gatewayLongitude := ac.GetAppSetting("GatewayLongitude")

		startupTime := time.Now().UnixNano() / int64(time.Millisecond)
		config.UUID = ac.GetAppSetting("GatewayId")
		config.Description = ac.GetAppSetting("GatewayDescription")
		config.Address = ac.GetAppSetting("GatewayHostname")
		config.Router = ac.GetAppSetting("GatewayRouter")
		config.RouterPort = ac.GetAppSetting("GatewayRouterPort")
		config.DeployNetwork = ac.GetAppSetting("GatewayDeployNetwork")
		config.Latitude, _ = strconv.ParseFloat(gatewayLatitude, 32)
		config.Longitude, _ = strconv.ParseFloat(gatewayLongitude, 32)
		config.AccessToken = ac.GetAppSetting("GatewayAccessToken")
		config.Username = ac.GetAppSetting("GatewayUsername")
		config.Platform = ac.GetAppSetting("GatewayPlatform")
		config.Createdts = startupTime
		config.Updatedts = startupTime
	} else {
		return nil, errors.New("no application-specific settings found")
	}

	return &config, nil
}
