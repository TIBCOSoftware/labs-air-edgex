package transforms

import (
	"encoding/json"
	"fmt"
	"sync"
)

// filterInfo - maintains device-resource information
type filterInfo struct {
	Device   string `json:"device"`
	Resource string `json:"resource"`
}

type filtersConfig struct {
	Filters []filterInfo `json:"filters"`
}

var filterInfoMap map[string]filterInfo

var mutex = &sync.Mutex{}

// InitializeFiltering -
func InitializeFiltering() {

	mutex.Lock()

	// Initialize Model map
	filterInfoMap = make(map[string]filterInfo)

	mutex.Unlock()

}

func addSampleFilterInfo() {

	// addNvidiaModelInfo("localtester", "image_recognition", "googlenet", "CameraPiHQ001", "onvif_snapshot")
}

func addFilterInfo(device string, resource string) {
	loggingClient.Info(fmt.Sprintf("Adding filter for: %s-%s", device, resource))

	fi := filterInfo{
		Device:   device,
		Resource: resource,
	}

	mapKey := fi.Device + "_" + fi.Resource

	mutex.Lock()
	filterInfoMap[mapKey] = fi
	mutex.Unlock()

}

// RegisterFilters - register filters
func RegisterFilters(config []byte) {

	InitializeFiltering()

	loggingClient.Debug(fmt.Sprintf("Processing config request preparse %v\n", config))

	fc := filtersConfig{}

	if err := json.Unmarshal([]byte(config), &fc); err != nil {
		fmt.Printf("Processing config request ERROR\n")
	}

	for _, filter := range fc.Filters {
		addFilterInfo(filter.Device, filter.Resource)
	}

}

// IsFiltered - indicates if the device-resource are registered for filtering and values
func IsFiltered(device string, resource string) bool {
	mapKey := device + "_" + resource
	mutex.Lock()
	_, ok := filterInfoMap[mapKey]
	mutex.Unlock()
	return ok
}
