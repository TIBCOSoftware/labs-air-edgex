package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"sync"
	"time"

	sdkModel "github.com/edgexfoundry/device-sdk-go/pkg/models"
	sdk "github.com/edgexfoundry/device-sdk-go/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/pkg/errors"
)

var once sync.Once
var lock sync.Mutex

var rtspClients map[string]*RtspClient

var driver *Driver

// Driver implements the sdkModel.ProtocolDriver interface for
// the device service
type Driver struct {
	lc       logger.LoggingClient
	asynchCh chan<- *sdkModel.AsyncValues
	config   *configuration
}

// NewProtocolDriver initializes the singleton Driver and
// returns it to the caller
func NewProtocolDriver() *Driver {
	once.Do(func() {
		driver = new(Driver)
		rtspClients = make(map[string]*RtspClient)
	})

	return driver
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (d *Driver) HandleReadCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []sdkModel.CommandRequest) ([]*sdkModel.CommandValue, error) {
	var responses = make([]*sdkModel.CommandValue, len(reqs))
	var resTime = time.Now().UnixNano() / int64(time.Millisecond)

	d.lc.Info(fmt.Sprintf("Handling read command for:  %s, Num Requests: %d", deviceName, len(reqs)))

	addr, err := d.addrFromProtocols(protocols)
	if err != nil {
		return responses, errors.Errorf("handleReadCommands: %v", err.Error())
	}

	// check for existence of both clients
	rtspClient, err := d.clientsFromAddr(addr, deviceName)
	if err != nil {
		return responses, errors.Errorf("handleReadCommands: %v", err.Error())
	}

	for i, req := range reqs {
		switch req.DeviceResourceName {
		// RTSP cases
		case "rtsp_device_information":
			d.lc.Info(fmt.Sprintf("Handling:  %s", "rtsp_device_information"))
			data, err := rtspClient.GetDeviceInformation()

			if err != nil {
				d.lc.Error(err.Error())
				return responses, err
			}

			cv := sdkModel.NewStringValue(reqs[i].DeviceResourceName, resTime, string(data))
			responses[i] = cv
		case "rtsp_profile_information":
			d.lc.Info(fmt.Sprintf("Handling:  %s", "rtsp_profile_information"))
			data, err := rtspClient.GetProfileInformation()

			if err != nil {
				d.lc.Error(err.Error())
				return responses, err
			}

			cv := sdkModel.NewStringValue(reqs[i].DeviceResourceName, 0, string(data))
			responses[i] = cv
		case "RtspDateTime":
			d.lc.Info(fmt.Sprintf("Handling:  %s", "RtspDateTime"))
			data, err := rtspClient.GetSystemDateAndTime()

			if err != nil {
				d.lc.Error(err.Error())
				return responses, err
			}

			cv := sdkModel.NewStringValue(reqs[i].DeviceResourceName, resTime, string(data))
			responses[i] = cv
		case "RtspHostname":
			d.lc.Info(fmt.Sprintf("Handling:  %s", "RtspHostname"))
			data, err := rtspClient.GetHostname()

			if err != nil {
				d.lc.Error(err.Error())
				return responses, err
			}

			cv := sdkModel.NewStringValue(reqs[i].DeviceResourceName, resTime, string(data))
			responses[i] = cv
		case "rtsp_dns":
			d.lc.Info(fmt.Sprintf("Handling:  %s", "rtsp_dns"))
			data, err := rtspClient.GetDNS()

			if err != nil {
				d.lc.Error(err.Error())
				return responses, err
			}

			cv := sdkModel.NewStringValue(reqs[i].DeviceResourceName, resTime, string(data))
			responses[i] = cv
		case "rtsp_network_interfaces":
			d.lc.Info(fmt.Sprintf("Handling:  %s", "rtsp_network_interfaces"))
			data, err := rtspClient.GetNetworkInterfaces()

			if err != nil {
				d.lc.Error(err.Error())
				return responses, err
			}

			cv := sdkModel.NewStringValue(reqs[i].DeviceResourceName, resTime, string(data))
			responses[i] = cv
		case "rtsp_network_protocols":
			d.lc.Info(fmt.Sprintf("Handling:  %s", "rtsp_network_protocols"))
			data, err := rtspClient.GetNetworkProtocols()

			if err != nil {
				d.lc.Error(err.Error())
				return responses, err
			}

			cv := sdkModel.NewStringValue(reqs[i].DeviceResourceName, resTime, string(data))
			responses[i] = cv
		case "rtsp_network_default_gateway":
			d.lc.Info(fmt.Sprintf("Handling:  %s", "rtsp_network_default_gateway"))
			data, err := rtspClient.GetNetworkDefaultGateway()

			if err != nil {
				d.lc.Error(err.Error())
				return responses, err
			}

			cv := sdkModel.NewStringValue(reqs[i].DeviceResourceName, resTime, string(data))
			responses[i] = cv
		case "rtsp_ntp":
			d.lc.Info(fmt.Sprintf("Handling:  %s", "rtsp_ntp"))
			data, err := rtspClient.GetNTP()

			if err != nil {
				d.lc.Error(err.Error())
				return responses, err
			}

			cv := sdkModel.NewStringValue(reqs[i].DeviceResourceName, resTime, string(data))
			responses[i] = cv
		case "rtsp_system_reboot":
			d.lc.Info(fmt.Sprintf("Handling:  %s", "rtsp_system_reboot"))
			data, err := rtspClient.Reboot()

			if err != nil {
				d.lc.Error(err.Error())
				return responses, err
			}

			cv := sdkModel.NewStringValue(reqs[i].DeviceResourceName, 0, string(data))
			responses[i] = cv
		case "rtsp_users":
			d.lc.Info(fmt.Sprintf("Handling:  %s", "rtsp_users"))
			data, err := rtspClient.GetUsers()

			if err != nil {
				d.lc.Error(err.Error())
				return responses, err
			}

			cv := sdkModel.NewStringValue(reqs[i].DeviceResourceName, resTime, string(data))
			responses[i] = cv
		case "rtsp_snapshot":
			d.lc.Info(fmt.Sprintf("Handling:  %s", "rtsp_snapshot"))
			data, err := rtspClient.GetSnapshot()

			saveImageToFile(data)

			if err != nil {
				d.lc.Error(err.Error())
				return responses, err
			}

			cv, err := sdkModel.NewBinaryValue(reqs[i].DeviceResourceName, resTime, data)

			if err != nil {
				err = errors.Wrap(err, "error creating binary CommandValue")
				d.lc.Error(err.Error())
				return responses, err
			}
			responses[i] = cv
		case "RtspStreamURI":
			d.lc.Info(fmt.Sprintf("Handling:  %s", "RtspStreamURI"))
			data, err := rtspClient.GetStreamURI()

			if err != nil {
				d.lc.Error(err.Error())
				return responses, err
			}

			cv := sdkModel.NewStringValue(reqs[i].DeviceResourceName, resTime, string(data))
			responses[i] = cv

		// camera specific cases
		default:
			d.lc.Info(fmt.Sprintf("Handling:  %s", "default"))
			err := errors.New("Non-RTSP command for camera")
			d.lc.Error(err.Error())
			return responses, err
		}
	}

	return responses, nil
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource (aka DeviceObject).
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (d *Driver) HandleWriteCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []sdkModel.CommandRequest, params []*sdkModel.CommandValue) error {
	addr, err := d.addrFromProtocols(protocols)
	if err != nil {
		return errors.Errorf("handleWriteCommands: %v", err.Error())
	}

	d.lc.Info(fmt.Sprintf("Handling write command for:  %s", deviceName))

	// check for existence of both clients
	rtspClient, err := d.clientsFromAddr(addr, deviceName)
	if err != nil {
		return errors.Errorf("handleWriteCommands: %v", err.Error())
	}

	for i, req := range reqs {
		switch req.DeviceResourceName {
		case "RtspUser":
			user := struct {
				Username  string
				Password  string
				UserLevel string
				Extension *string
			}{}

			err := structFromParam(params[i], &user)
			if err != nil {
				d.lc.Error(err.Error())
				return err
			}

			err = rtspClient.CreateUser("")
			if err != nil {
				d.lc.Error(fmt.Sprintf("handleWriteCommands error: %v", err.Error()))
				return err
			}

		case "RtspReboot":
			shouldReboot, err := params[i].BoolValue()
			if err != nil {
				err := errors.New("non-binary value passed to RtspReboot command")
				d.lc.Error(err.Error())
				return err
			}
			if !shouldReboot {
				continue
			}

			_, err = rtspClient.Reboot()
			if err != nil {
				return err
			}

		case "RtspHostname":
			hostname, err := params[i].StringValue()
			if err != nil {
				err := errors.New("non-string value passed to RtspHostname command")
				d.lc.Error(err.Error())
				return err
			}

			err = rtspClient.SetHostname(hostname)
			if err != nil {
				d.lc.Error(err.Error())
				return err
			}

		case "RtspHostnameFromDHCP":
			err := rtspClient.SetHostnameFromDHCP()
			if err != nil {
				d.lc.Error(err.Error())
				return err
			}

		case "RtspDateTime":
			dateTime := struct {
				Year   int
				Month  int
				Day    int
				Hour   int
				Minute int
				Second int
			}{}

			err := structFromParam(params[i], &dateTime)
			if err != nil {
				d.lc.Error(err.Error())
				return err
			}

			t := time.Date(dateTime.Year, time.Month(dateTime.Month), dateTime.Day, dateTime.Hour, dateTime.Minute, dateTime.Second, 0, time.UTC)
			err = rtspClient.SetSystemDateAndTime(t)
			if err != nil {
				d.lc.Error(err.Error())
				return err
			}

		default:
			err := errors.New("non-rtsp command for camera")
			d.lc.Error(err.Error())
			return err
		}

	}

	return nil
}

type stringer interface {
	StringValue() (string, error)
}

func structFromParam(s stringer, v interface{}) error {
	str, err := s.StringValue()
	if err != nil {
		return errors.Errorf("RtspUser CommandValue missing string value")
	}
	err = json.Unmarshal([]byte(str), v)
	if err != nil {
		return errors.Errorf("error unmarshaling string: %v", err.Error())
	}
	return nil
}

// DisconnectDevice handles protocol-specific cleanup when a device
// is removed.
func (d *Driver) DisconnectDevice(deviceName string, protocols map[string]contract.ProtocolProperties) error {
	addr, err := d.addrFromProtocols(protocols)
	if err != nil {
		return errors.Errorf("no address found for device: %v", err.Error())
	}

	shutdownClient(addr)
	shutdownRtspClient(addr)
	return nil
}

// Initialize performs protocol-specific initialization for the device
// service.
func (d *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *sdkModel.AsyncValues, deviceCh chan<- []sdkModel.DiscoveredDevice) error {
	d.lc = lc
	d.asynchCh = asyncCh

	d.lc.Info(fmt.Sprintf("Initializing device:  %s", "camera"))

	config, err := loadCameraConfig()

	if err != nil {
		panic(fmt.Errorf("load camera configuration failed: %d", err))
	}
	d.config = config

	for _, dev := range sdk.RunningService().Devices() {
		d.lc.Info(fmt.Sprintf("Initializing RSTP clients:  %s", "camera"))
		initializeRtspClient(dev, config.Camera.User, config.Camera.Password, config.Camera.AuthMethod)
	}

	return nil
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (d *Driver) Stop(force bool) error {

	close(d.asynchCh)

	return nil
}

// AddDevice is a callback function that is invoked
// when a new Device associated with this Device Service is added
func (d *Driver) AddDevice(deviceName string, protocols map[string]contract.ProtocolProperties, adminState contract.AdminState) error {

	d.lc.Info(fmt.Sprintf("Adding device:  %s", deviceName))

	addr, err := d.addrFromProtocols(protocols)
	if err != nil {
		err = errors.Errorf("error adding device: %v", err.Error())
		d.lc.Error(err.Error())
		return err
	}

	_, err = d.clientsFromAddr(addr, deviceName)
	if err != nil {
		err = errors.Errorf("error adding device: %v", err.Error())
		d.lc.Error(err.Error())
		return err
	}
	return nil
}

// UpdateDevice is a callback function that is invoked
// when a Device associated with this Device Service is updated
func (d *Driver) UpdateDevice(deviceName string, protocols map[string]contract.ProtocolProperties, adminState contract.AdminState) error {
	return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (d *Driver) RemoveDevice(deviceName string, protocols map[string]contract.ProtocolProperties) error {
	addr, err := d.addrFromProtocols(protocols)
	if err != nil {
		return errors.Errorf("no address found for device: %v", err.Error())
	}

	shutdownClient(addr)
	shutdownRtspClient(addr)
	return nil
}

func getRtspClient(addr string) (*RtspClient, bool) {
	lock.Lock()
	c, ok := rtspClients[addr]
	lock.Unlock()
	return c, ok
}

func initializeRtspClient(device contract.Device, user string, password string, authMethod string) *RtspClient {
	addr := device.Protocols["RTSP"]["Address"]
	c := NewRtspClient(addr, user, password, authMethod, driver.lc)
	lock.Lock()
	rtspClients[addr] = c
	lock.Unlock()
	return c
}

func shutdownRtspClient(addr string) {
	// nothing much to do here at the moment
	lock.Lock()
	delete(rtspClients, addr)
	lock.Unlock()
}

func shutdownClient(addr string) {
	lock.Lock()

	lock.Unlock()
}

func in(needle string, haystack []string) bool {
	for _, e := range haystack {
		if needle == e {
			return true
		}
	}
	return false
}

func (d *Driver) addrFromProtocols(protocols map[string]contract.ProtocolProperties) (string, error) {

	if _, ok := protocols["RTSP"]; !ok {
		d.lc.Error("No RTSP address found for device. Check configuration file.")
		return "", fmt.Errorf("no HTTP address in protocols map")
	}

	var addr string
	addr, ok := protocols["RTSP"]["Address"]
	if !ok {
		d.lc.Error("No RTSP address found for device. Check configuration file.")
		return "", fmt.Errorf("no RTSP address in protocols map")
	}
	return addr, nil

}

func (d *Driver) clientsFromAddr(addr string, deviceName string) (*RtspClient, error) {
	rtspClient, ok := getRtspClient(addr)

	if !ok {
		dev, err := sdk.RunningService().GetDeviceByName(deviceName)
		if err != nil {
			err = fmt.Errorf("device not found: %s", deviceName)
			d.lc.Error(err.Error())

			return nil, err
		}

		rtspClient = initializeRtspClient(dev, d.config.Camera.User, d.config.Camera.Password, d.config.Camera.AuthMethod)
	}

	return rtspClient, nil
}

func saveImageToFile(imgByte []byte) {

	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Fatalln(err)
	}

	out, _ := os.Create("./rtspsnapshot.jpeg")
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 90

	err = jpeg.Encode(out, img, &opts)
	//jpeg.Encode(out, img, nil)
	if err != nil {
		log.Println(err)
	}

}
