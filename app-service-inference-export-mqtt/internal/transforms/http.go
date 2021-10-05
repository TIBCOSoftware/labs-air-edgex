//
// Copyright (c) 2019 Intel Corporation
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

package transforms

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/edgexfoundry/app-functions-sdk-go/pkg/util"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

// HTTPSender ...
type HTTPSender struct {
	URL            string
	MimeType       string
	PersistOnError bool
}

// NewHTTPSender creates, initializes and returns a new instance of HTTPSender
func NewHTTPSender(url string, mimeType string, persistOnError bool) *HTTPSender {
	return &HTTPSender{
		URL:            url,
		MimeType:       mimeType,
		PersistOnError: persistOnError,
	}
}

// HTTPPost will send data from the previous function to the specified Endpoint via http POST.
// If no previous function exists, then the event that triggered the pipeline will be used.
// An empty string for the mimetype will default to application/json.
func (sender HTTPSender) HTTPPost(logging logger.LoggingClient, msg string) (bool, interface{}) {

	if sender.MimeType == "" {
		sender.MimeType = "application/json"
	}
	exportData, err := util.CoerceType(msg)
	if err != nil {
		return false, err
	}

	response, err := http.Post(sender.URL, sender.MimeType, bytes.NewReader(exportData))
	if err != nil {
		logging.Debug(fmt.Sprintf("Error posting: %s", err))
		return false, err
	}

	defer response.Body.Close()
	logging.Debug(fmt.Sprintf("Response: %s", response.Status))

	return true, nil

}

// Predict calls model REST API
func (sender HTTPSender) Predict(logging logger.LoggingClient, modelURL string, jsonPayload string) (bool, interface{}) {

	if sender.MimeType == "" {
		sender.MimeType = "application/json"
	}

	requestBody, err := util.CoerceType(jsonPayload)
	if err != nil {
		return false, err
	}

	logging.Info(fmt.Sprintf("ModelURL: %s", modelURL))
	logging.Debug(fmt.Sprintf("Request: %v", requestBody))

	response, err := http.Post(modelURL, sender.MimeType, bytes.NewReader(requestBody))
	if err != nil {
		logging.Debug(fmt.Sprintf("Error posting: %s", err))
		return false, err
	}

	defer response.Body.Close()
	logging.Info(fmt.Sprintf("Response Status: %s", response.Status))

	// read response body
	data, _ := ioutil.ReadAll(response.Body)

	// logging.Info(fmt.Sprintf("Response Body: %s", data))

	return true, data

}
