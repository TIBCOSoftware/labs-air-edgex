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

const user = "User"
const password = "Password"
const authMethod = "AuthMethod"

// loadCameraConfig loads the camera configuration
func loadCameraConfig() (*configuration, error) {
	config := new(configuration)
	if val, ok := sdk.DriverConfigs()[user]; ok {
		config.Camera.User = val
	} else {
		return config, fmt.Errorf("driver config undefined: %s", user)
	}
	if val, ok := sdk.DriverConfigs()[password]; ok {
		config.Camera.Password = val
	} else {
		return config, fmt.Errorf("driver config undefined: %s", password)
	}
	if val, ok := sdk.DriverConfigs()[authMethod]; ok {
		config.Camera.AuthMethod = val
	} else {
		return config, fmt.Errorf("driver config undefined: %s", authMethod)
	}

	return config, nil
}
