package driver

import (
	"fmt"

	sdk "github.com/edgexfoundry/device-sdk-go/pkg/service"
)

type configuration struct {
	Camera cameraInfo
}

type cameraInfo struct {
	User       string
	Password   string
	AuthMethod string
}

const (
	// USER - user
	USER = "User"
	// PASSWORD - password
	PASSWORD = "Password"
	// AUTHMETHOD - authentication method
	AUTHMETHOD = "AuthMethod"
)

// loadCameraConfig loads the camera configuration
func loadCameraConfig() (*configuration, error) {
	config := new(configuration)
	if val, ok := sdk.DriverConfigs()[USER]; ok {
		config.Camera.User = val
	} else {
		return config, fmt.Errorf("driver config undefined: %s", USER)
	}
	if val, ok := sdk.DriverConfigs()[PASSWORD]; ok {
		config.Camera.Password = val
	} else {
		return config, fmt.Errorf("driver config undefined: %s", PASSWORD)
	}
	if val, ok := sdk.DriverConfigs()[AUTHMETHOD]; ok {
		config.Camera.AuthMethod = val
	} else {
		return config, fmt.Errorf("driver config undefined: %s", AUTHMETHOD)
	}

	return config, nil
}
