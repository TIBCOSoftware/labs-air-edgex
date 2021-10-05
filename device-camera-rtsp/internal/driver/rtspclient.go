package driver

import (
	"fmt"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"gocv.io/x/gocv"
)

// RtspClient manages the state required to issue ONVIF requests to a camera
type RtspClient struct {
	ipAddress  string
	user       string
	password   string
	cameraAuth string
	url        string
	lc         logger.LoggingClient
}

func getRstpImage(url string) ([]byte, error) {

	webcam, err := gocv.OpenVideoCapture(url)
	defer webcam.Close()

	if err != nil {
		fmt.Printf("Error opening RTSP device\n")
		return nil, err
	}

	img := gocv.NewMat()
	defer img.Close()

	if ok := webcam.Read(&img); !ok {
		fmt.Printf("RTSP Device closed \n")
		return nil, fmt.Errorf("Error reading webcam")
	}

	if img.Empty() {
		return nil, fmt.Errorf("Error reading image")
	}

	outbytes, err := gocv.IMEncode(".jpg", img)
	if err != nil {
		return nil, fmt.Errorf("Error in IMEncode")
	}

	return outbytes, nil

}

// NewRtspClient returns an RtspClient for a single camera
func NewRtspClient(ipAddress string, user string, password string, cameraAuth string, lc logger.LoggingClient) *RtspClient {

	rtspURL := "rtsp://" + user + ":" + password + "@" + ipAddress

	c := RtspClient{
		ipAddress:  ipAddress,
		user:       user,
		password:   password,
		cameraAuth: cameraAuth,
		url:        rtspURL,
		lc:         lc,
	}

	lc.Info(fmt.Sprintf("NewRtspClient with ip: %s", ipAddress))

	rtspSrc := "rtsp://" + user + ":" + password + "@" + ipAddress

	lc.Info(fmt.Sprintf("RtspClient url: %s", rtspSrc))

	return &c
}

// GetDeviceInformation makes an ONVIF GetDeviceInformation request to the camera
func (c *RtspClient) GetDeviceInformation() (string, error) {

	return "", fmt.Errorf("Operation not supported")
}

// GetProfileInformation makes an ONVIF GetProfiles request to the camera
func (c *RtspClient) GetProfileInformation() (string, error) {

	return "", fmt.Errorf("Operation not supported")
}

// GetStreamURI returns the RTSP URI for the first media profile returned by the camera
func (c *RtspClient) GetStreamURI() (string, error) {

	return "", fmt.Errorf("Operation not supported")
}

// GetSnapshot returns a snapshot from the camera as a slice of bytes
func (c *RtspClient) GetSnapshot() ([]byte, error) {

	c.lc.Info(fmt.Sprintf("GetSnapshot for url: %s", c.url))

	return getRstpImage(c.url)
}

// GetSystemDateAndTime returns the current date and time as reported by the ONVIF GetSystemDateAndTime command
func (c *RtspClient) GetSystemDateAndTime() (string, error) {

	return "", fmt.Errorf("Operation not supported")
}

// GetHostname returns the hostname reported by the device via the ONVIF GetHostname command
func (c *RtspClient) GetHostname() (string, error) {

	return "", fmt.Errorf("Operation not supported")
}

// SetHostname requests a change to the camera's hostname via the ONFVIF SetHostname command
func (c *RtspClient) SetHostname(name string) error {

	return fmt.Errorf("Operation not supported")
}

// SetHostnameFromDHCP requests the camera to base its hostname from DHCP
func (c *RtspClient) SetHostnameFromDHCP() error {

	return fmt.Errorf("Operation not supported")
}

// SetSystemDateAndTime changes the camera's system time via the SetSystemDateAndTime ONVIF command
func (c *RtspClient) SetSystemDateAndTime(datetime time.Time) error {

	return fmt.Errorf("Operation not supported")
}

// GetDNS returns the DNS settings as reported by the ONVIF GetDNS command
func (c *RtspClient) GetDNS() (string, error) {

	return "", fmt.Errorf("Operation not supported")
}

// GetNetworkInterfaces returns the results of the ONVIF GetNetworkInterfaces command
func (c *RtspClient) GetNetworkInterfaces() (string, error) {

	return "", fmt.Errorf("Operation not supported")
}

// GetNetworkProtocols returns the resutls of the ONVIF GetNetworkProtocols command
func (c *RtspClient) GetNetworkProtocols() (string, error) {

	return "", fmt.Errorf("Operation not supported")
}

// GetNetworkDefaultGateway returns the results of the ONVIF GetNetworkDefaultGateway command
func (c *RtspClient) GetNetworkDefaultGateway() (string, error) {

	return "", fmt.Errorf("Operation not supported")
}

// GetNTP returns the results of the ONVIF GetNTP command
func (c *RtspClient) GetNTP() (string, error) {

	return "", fmt.Errorf("Operation not supported")
}

// Reboot requests a device system reboot via ONVIF
func (c *RtspClient) Reboot() (string, error) {

	return "", fmt.Errorf("Operation not supported")
}

// GetUsers requests the Users associated with the device via ONVIF
func (c *RtspClient) GetUsers() (string, error) {

	return "", fmt.Errorf("Operation not supported")
}

// CreateUser creates a new ONVIF User for the device
func (c *RtspClient) CreateUser(user string) error {

	return fmt.Errorf("Operation not supported")
}
