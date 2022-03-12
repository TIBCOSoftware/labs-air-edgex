// TODO: Change Copyright to your company if open sourcing or remove header
//
// Copyright (c) 2021 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"os"
	"strconv"
	"time"

	"github.com/TIBCOSoftware/labs-air/edgexfoundry/app-service-metadata/internal/config"
	"github.com/TIBCOSoftware/labs-air/edgexfoundry/app-service-metadata/internal/functions"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
)

const (
	serviceKey = "app-service-metadata"
)

// TODO: Define your app's struct
type myApp struct {
	service interfaces.ApplicationService
	lc      logger.LoggingClient
}

func main() {

	app := myApp{}

	var ok bool
	var err error
	app.service, ok = pkg.NewAppService(serviceKey)
	if !ok {
		os.Exit(1)
	}

	app.lc = app.service.LoggingClient()

	// Get the application's specific configuration settings.
	appSettings := app.service.ApplicationSettings()
	appConfig := config.NewAppConfig(app.lc, appSettings)

	// Create Export functions
	export := functions.NewExport(app.lc, appConfig, app.service.DeviceClient(), app.service.DeviceProfileClient())

	// This is out pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	err = app.service.SetFunctionsPipeline(
		export.ProcessEvent)
	if err != nil {
		app.lc.Errorf("SetFunctionsPipeline returned error: %s", err.Error())
		os.Exit(1)
	}

	// Set ticker to trigger metadata publishing
	metadataPublishInterval := 30
	metadataPublishIntervalStr := appConfig.GetAppSetting("MetadataPublishIntervalSecs")
	metadataPublishInterval64, _ := strconv.ParseInt(metadataPublishIntervalStr, 10, 32)
	metadataPublishInterval = int(metadataPublishInterval64)
	ticker := time.NewTicker(time.Duration(metadataPublishInterval) * time.Second)
	go func() {
		for ; true; <-ticker.C {
			export.PublishMetadata()
		}
	}()

	if err := app.service.MakeItRun(); err != nil {
		app.lc.Errorf("MakeItRun returned error: %s", err.Error())
		os.Exit(1)
	}

	// Do any required cleanup here, if needed
	ticker.Stop()
	os.Exit(0)
}
