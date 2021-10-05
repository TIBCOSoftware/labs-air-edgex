package transforms

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

// ModelResponse - struct used to return model inference responses
type ModelResponse struct {
	Predict string  `json:"predict"`
	Score   float32 `json:"score"`
}

type resnetStruct struct {
	B64 string `json:"b64"`
}

type resnetPredictStruct struct {
	Instances []resnetStruct `json:"instances"`
}

type rcnnresnetPredictStruct struct {
	SignatureName string      `json:"signature_name"`
	Instances     [][][][]int `json:"instances"`
}

type imagenetPredictStruct struct {
	Network string `json:"network"`
	Image   string `json:"image"`
}

type patternPredictStruct struct {
	Data string `json:"Data"`
}

type rcnnresnetPredictionStruct struct {
	NumDetections    float32     `json:"num_detections"`
	DetectionBoxes   [][]float32 `json:"detection_boxes"`
	DetectionScores  []float32   `json:"detection_scores"`
	DetectionClasses []float32   `json:"detection_classes"`
}

type rcnnresnetPredictResponseStruct struct {
	Predictions []rcnnresnetPredictionStruct `json:"predictions"`
}

type imagenetPredictResponseStruct struct {
	ClassName  string  `json:"classname"`
	ClassIndex int     `json:"classindex"`
	Confidence float32 `json:"confidence"`
}

type patternStruct struct {
	LotWaferID string  `json:"LotWaferID"`
	Pattern    string  `json:"Pattern"`
	Max        float32 `json:"Max"`
}

type patternRecognitionPredictResponseStruct struct {
	Prediction patternStruct `json:"prediction"`
}

type anomalyDetectionReadingStruct struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

type anomalyDetectionStruct struct {
	Series          []anomalyDetectionReadingStruct `json:"series"`
	Granularity     string                          `json:"granularity"`
	MaxAnomalyRatio float32                         `json:"maxAnomalyRatio"`
	Sensitivity     int                             `json:"sensitivity"`
}

type anomalyDetectionResponseStruct struct {
	Period            int     `json:"series"`
	SuggestedWindow   int     `json:"suggestedWindow"`
	ExpectedValue     float64 `json:"expectedValue"`
	UpperMargin       float64 `json:"upperMargin"`
	LowerMargin       float64 `json:"lowerMargin"`
	IsAnomaly         bool    `json:"isAnomaly"`
	IsNegativeAnomaly bool    `json:"isNegativeAnomaly"`
	IsPositiveAnomaly bool    `json:"isPositiveAnomaly"`
}

var cocoCategories = map[int]string{
	0:  "background",
	1:  "person",
	2:  "bicycle",
	3:  "car",
	4:  "motorcycle",
	5:  "airplane",
	6:  "bus",
	7:  "train",
	8:  "truck",
	9:  "boat",
	10: "traffic light",
	11: "fire hydrant",
	12: "12",
	13: "stop sign",
	14: "parking meter",
	15: "bench",
	16: "bird",
	17: "cat",
	18: "dog",
	19: "horse",
	20: "sheep",
	21: "cow",
	22: "elephant",
	23: "bear",
	24: "zebra",
	25: "giraffe",
	26: "26",
	27: "backpack",
	28: "umbrella",
	29: "29",
	30: "30",
	31: "handbag",
	32: "tie",
	33: "suitcase",
	34: "frisbee",
	35: "skis",
	36: "snowboard",
	37: "sports ball",
	38: "kite",
	39: "baseball bat",
	40: "baseball glove",
	41: "skateboard",
	42: "surfboard",
	43: "tennis racket",
	44: "bottle",
	45: "45",
	46: "wine glass",
	47: "cup",
	48: "fork",
	49: "knife",
	50: "spoon",
	51: "bowl",
	52: "banana",
	53: "apple",
	54: "sandwich",
	55: "orange",
	56: "broccoli",
	57: "carrot",
	58: "hot dog",
	59: "pizza",
	60: "donut",
	61: "cake",
	62: "chair",
	63: "couch",
	64: "potted plant",
	65: "bed",
	66: "66",
	67: "dining table",
	68: "68",
	69: "69",
	70: "toilet",
	71: "71",
	72: "tv",
	73: "laptop",
	74: "mouse",
	75: "remote",
	76: "keyboard",
	77: "cell phone",
	78: "microwave",
	79: "oven",
	80: "toaster",
	81: "sink",
	82: "refrigerator",
	83: "83",
	84: "book",
	85: "clock",
	86: "vase",
	87: "scissors",
	88: "teddy bear",
	89: "hair drier",
	90: "toothbrush",
}

var adReadingsQueue = []float64{
	256.651276, 267.87323, 270.50769, 257.070831, 248.238388, 243.320313, 236.744781, 229.485672, 226.458328, 227.597641,
	229.135727, 230.885406, 225.159164, 214.457916, 203.762939, 191.954147, 179.334976, 172.94339, 169.314377, 171.65538,
	168.799347, 166.127533, 159.339645, 151.679764, 142.931793, 137.939377, 131.901886, 125.198395, 114.879669, 107.6092,
	104.036125, 100.041542, 96.2975082, 94.72258, 94.8925476, 96.5237274, 96.5348282, 96.9986877, 97.9477692, 99.0746536,
	98.6401367, 99.8225861, 100.742188, 100.083008, 99.4319611, 100.332344, 101.117027,
}

var adTimestampQueue = []string{
	"2016-01-01T00:00:00Z", "2016-02-01T00:00:00Z", "2016-03-01T00:00:00Z", "2016-04-01T00:00:00Z",
	"2016-05-01T00:00:00Z", "2016-06-01T00:00:00Z", "2016-07-01T00:00:00Z", "2016-08-01T00:00:00Z",
	"2016-09-01T00:00:00Z", "2016-10-01T00:00:00Z", "2016-11-01T00:00:00Z", "2016-12-01T00:00:00Z",
	"2017-01-01T00:00:00Z", "2017-02-01T00:00:00Z", "2017-03-01T00:00:00Z", "2017-04-01T00:00:00Z",
	"2017-05-01T00:00:00Z", "2017-06-01T00:00:00Z", "2017-07-01T00:00:00Z", "2017-08-01T00:00:00Z",
	"2017-09-01T00:00:00Z", "2017-10-01T00:00:00Z", "2017-11-01T00:00:00Z", "2017-12-01T00:00:00Z",
	"2018-01-01T00:00:00Z", "2018-02-01T00:00:00Z", "2018-03-01T00:00:00Z", "2018-04-01T00:00:00Z",
	"2018-05-01T00:00:00Z", "2018-06-01T00:00:00Z", "2018-07-01T00:00:00Z", "2018-08-01T00:00:00Z",
	"2018-09-01T00:00:00Z", "2018-10-01T00:00:00Z", "2018-11-01T00:00:00Z", "2018-12-01T00:00:00Z",
	"2019-01-01T00:00:00Z", "2019-02-01T00:00:00Z", "2019-03-01T00:00:00Z", "2019-04-01T00:00:00Z",
	"2019-05-01T00:00:00Z", "2019-06-01T00:00:00Z", "2019-07-01T00:00:00Z", "2019-08-01T00:00:00Z",
	"2019-09-01T00:00:00Z", "2019-10-01T00:00:00Z", "2019-11-01T00:00:00Z", "2019-12-01T00:00:00Z",
}

// ModelInfo - maintains device-model information
type modelInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Version     string `json:"version"`
	Device      string `json:"device"`
	Resource    string `json:"resource"`
	ModelType   string `json:"modelType"`
	Model       string `json:"model"`
	Server      string `json:"server"`
}

var modelInfoMap map[string]modelInfo

var httpSender *HTTPSender

var tfInferenceServiceURL string
var msInferenceServiceURL string
var nvidiaInferenceServiceURL string
var tibcoInferenceServiceURL string

// InitializeModeling - initialize sender
func InitializeModeling(settings map[string]string) {
	// Initialize Inference Service URLs

	tfInferenceServiceURL = GetAppSetting(settings, "TFInferenceServiceURL")
	msInferenceServiceURL = GetAppSetting(settings, "MSInferenceServiceURL")
	nvidiaInferenceServiceURL = GetAppSetting(settings, "NvidiaInferenceServiceURL")
	tibcoInferenceServiceURL = GetAppSetting(settings, "TIBCOInferenceServiceURL")

	// Initialize Model map
	modelInfoMap = make(map[string]modelInfo)

	// Create HTTP Sender
	httpSender = NewHTTPSender("", "application/json", false)

	// addSampleModelInfo()
}

func addSampleModelInfo() {

	addNvidiaModelInfo("localtester", "image_recognition", "googlenet", "CameraPiHQ001", "onvif_snapshot")
}

func addTIBCOModelInfo(name string, modelType string, model string, device string, resource string) {
	loggingClient.Info(fmt.Sprintf("Adding registration for TIBCO Model: %s-%s %s-%s", modelType, model, device, resource))

	mi := modelInfo{
		Name:      name,
		URL:       tibcoInferenceServiceURL + "/v1/" + modelType,
		Device:    device,
		Resource:  resource,
		ModelType: modelType,
		Model:     model,
		Server:    "nvidia",
	}

	mapKey := mi.Device + "_" + mi.Resource

	modelInfoMap[mapKey] = mi

}

func addNvidiaModelInfo(name string, modelType string, model string, device string, resource string) {
	loggingClient.Info(fmt.Sprintf("Adding registration for Nvidia Model: %s-%s %s-%s", modelType, model, device, resource))

	mi := modelInfo{
		Name:      name,
		URL:       nvidiaInferenceServiceURL + "/v1/" + modelType,
		Device:    device,
		Resource:  resource,
		ModelType: modelType,
		Model:     model,
		Server:    "nvidia",
	}

	mapKey := mi.Device + "_" + mi.Resource

	modelInfoMap[mapKey] = mi

}

func addMSModelInfo(name string, modelType string, model string, device string, resource string) {
	loggingClient.Info(fmt.Sprintf("Adding registration for MS Model: %s %s", device, resource))

	mi := modelInfo{
		Name:      name,
		URL:       msInferenceServiceURL + "/anomalydetector/v1.0/timeseries/last/detect",
		Device:    device,
		Resource:  resource,
		ModelType: modelType,
		Model:     model,
		Server:    "ms",
	}

	mapKey := mi.Device + "_" + mi.Resource

	modelInfoMap[mapKey] = mi

}

func addTFServingModelInfo(name string, modelType string, model string, device string, resource string) {
	loggingClient.Info(fmt.Sprintf("Adding registration for TFServing Model for %s: %s", device, resource))

	mi := modelInfo{
		Name:      name,
		URL:       tfInferenceServiceURL + "/v1/models/" + model + ":predict",
		Device:    device,
		Resource:  resource,
		ModelType: modelType,
		Model:     model,
		Server:    "tfserving",
	}

	mapKey := mi.Device + "_" + mi.Resource

	modelInfoMap[mapKey] = mi

}

// RegisterModel - register model
func RegisterModel(config []byte) {

	mi := modelInfo{}

	if err := json.Unmarshal([]byte(config), &mi); err != nil {
		fmt.Printf("Processing config request ERROR\n")
	}

	modelDetails := strings.Split(mi.Model, "|")

	loggingClient.Info(fmt.Sprintf("Register Model: %v", mi))
	loggingClient.Info(fmt.Sprintf("Register Model Details: %v", modelDetails))

	if modelDetails[0] == "tibco" {
		addTIBCOModelInfo(mi.Name, modelDetails[1], modelDetails[2], mi.Device, mi.Resource)
	} else if modelDetails[0] == "nvidia" {
		addNvidiaModelInfo(mi.Name, modelDetails[1], modelDetails[2], mi.Device, mi.Resource)
	} else if modelDetails[0] == "tfserving" {
		addTFServingModelInfo(mi.Name, modelDetails[1], modelDetails[2], mi.Device, mi.Resource)
	} else if modelDetails[0] == "ms" {
		addMSModelInfo(mi.Name, modelDetails[1], modelDetails[2], mi.Device, mi.Resource)
	} else {
		fmt.Printf("Processing Register Model ERROR\n")
	}

}

// UnregisterModel - register model
func UnregisterModel(config []byte) {

	mi := modelInfo{}

	if err := json.Unmarshal([]byte(config), &mi); err != nil {
		fmt.Printf("Processing config request ERROR\n")
	}

	mapKey := mi.Device + "_" + mi.Resource

	delete(modelInfoMap, mapKey)

}

// IsInferable - indicates if the device-resource are registered for modeling and values
// can be infered
func IsInferable(device string, resource string) bool {
	mapKey := device + "_" + resource
	_, ok := modelInfoMap[mapKey]
	return ok
}

func getModelInfo(device string, resource string) modelInfo {
	mapKey := device + "_" + resource
	return modelInfoMap[mapKey]
}

// BuildRequestFromBinary - build request to be passed to model's REST API
func buildRequestFromBinary(mi modelInfo, binaryValue []byte) ([]byte, error) {

	switch {
	case mi.Server == "nvidia":
		fmt.Println("Building nvidia binary request")
		return buildNvidiaRequestFromBinary(mi, binaryValue)
	case mi.Server == "tfserving":
		fmt.Println("Building tfserving binary request")
		if mi.Model == "renet" {
			fmt.Println("Building resnet request")
			return buildResnetRequestFromBinary(binaryValue)
		} else if mi.Model == "rcnnresnet" {
			fmt.Println("Building rcnnresnet request.")
			return buildRcnnResnetRequestFromBinary(binaryValue)
		} else {
			fmt.Println("buildRequestFromBinary - Invalid model")
		}
	default:
		fmt.Println("buildRequestFromBinary - Invalid vendor")
	}

	return nil, errors.New("Not a supported model")
}

// BuildRequestFromPattern - build request to be passed to model's REST API
func buildRequestFromPattern(mi modelInfo, pattern string) ([]byte, error) {

	// Pattern string needs to be decoded before sending to model prediction
	decodedPattern, _ := base64.StdEncoding.DecodeString(pattern)

	jsonPatternPredictStruc := &patternPredictStruct{
		Data: string(decodedPattern),
	}

	encjson, error := json.Marshal(jsonPatternPredictStruc)

	return encjson, error
}

func buildResnetRequestFromBinary(binaryValue []byte) ([]byte, error) {

	encimg := base64.StdEncoding.EncodeToString(binaryValue)

	jsonResnetStruct := resnetStruct{
		B64: encimg,
	}

	jsonResnetPredictStruc := &resnetPredictStruct{
		Instances: []resnetStruct{jsonResnetStruct},
	}

	encjson, error := json.Marshal(jsonResnetPredictStruc)

	return encjson, error
}

func buildRcnnResnetRequestFromBinary(binaryValue []byte) ([]byte, error) {

	imageData, imageType, err := image.Decode(bytes.NewReader(binaryValue))

	if err == nil {
		fmt.Printf("Received Image Type: %s, Image Size: %s, Color in middle: %v\n",
			imageType, imageData.Bounds().Size().String(),
			imageData.At(imageData.Bounds().Size().X/2, imageData.Bounds().Size().Y/2))
	}

	r := imageData.Bounds()
	width := r.Max.X
	height := r.Max.Y

	data := make([][][][]int, 1)

	data[0] = make([][][]int, height)

	for y := 0; y < height; y++ {
		data[0][y] = make([][]int, width)
		for x := 0; x < width; x++ {

			p := imageData.At(x, y)
			c := color.NRGBAModel.Convert(p).(color.NRGBA)

			data[0][y][x] = make([]int, 3)
			data[0][y][x][0] = int(c.R)
			data[0][y][x][1] = int(c.G)
			data[0][y][x][2] = int(c.B)
		}
	}

	jsonRcnnResnetPredictStruc := &rcnnresnetPredictStruct{
		SignatureName: "serving_default",
		Instances:     data,
	}

	encjson, error := json.Marshal(jsonRcnnResnetPredictStruc)

	// fmt.Println(string(encjson))

	return encjson, error
}

func buildNvidiaRequestFromBinary(mi modelInfo, binaryValue []byte) ([]byte, error) {

	imageData := base64.StdEncoding.EncodeToString(binaryValue)

	jsonImagenetPredictStruc := &imagenetPredictStruct{
		Network: mi.Model,
		Image:   imageData,
	}

	encjson, error := json.Marshal(jsonImagenetPredictStruc)

	// fmt.Println(string(encjson))

	return encjson, error
}

func buildAnomalyDetectionRequest(value float64) ([]byte, error) {

	// Append value to queue
	adReadingsQueue = append(adReadingsQueue, value)

	// Create readings
	readings := make([]anomalyDetectionReadingStruct, 48)

	for i, v := range adTimestampQueue {
		readings[i].Timestamp = v
		readings[i].Value = adReadingsQueue[i]
	}

	jsondat := &anomalyDetectionStruct{
		Series:          readings,
		Granularity:     "monthly",
		MaxAnomalyRatio: 0.25,
		Sensitivity:     95,
	}

	encjson, error := json.Marshal(jsondat)

	// fmt.Println(string(encjson))

	// Dequeue first element
	_, adReadingsQueue = adReadingsQueue[0], adReadingsQueue[1:]

	return encjson, error
}

// ParseResponse - parses the result from a model into a common interface
func parseResponse(mi modelInfo, response []byte) (*ModelResponse, error) {

	switch {
	case mi.ModelType == "pattern_recognition":
		fmt.Println("Building pattern recognition response")
		return parsePatternRecognitionResponse(response)
	case mi.Server == "nvidia":
		fmt.Println("Building nvidia response")
		return parseImagenetResponse(response)
	case mi.Server == "tfserving":
		if mi.Model == "renet" {
			fmt.Println("Building resnet response")
		} else if mi.Model == "rcnnresnet" {
			fmt.Println("Building rcnnresnet response")
			return parseRcnnResnetResponse(response)
		} else {
			fmt.Println("Parse response - Invalid model")
		}
	case mi.Server == "ms":
		if mi.Model == "AnomalyDetection" {
			fmt.Println("Building anomaly detection response")
			return parseAnomalyDetectionResponse(response)
		}
	default:
		fmt.Println("Parse response - No model")
	}

	return nil, errors.New("Not a supported model")
}

func parseRcnnResnetResponse(response []byte) (*ModelResponse, error) {

	res := rcnnresnetPredictResponseStruct{}

	err := json.Unmarshal(response, &res)

	// fmt.Printf("Unmarshalled Response: %+v\n ", res)

	// fmt.Printf("Response NumDetections: %f  Classes: %v ", res.Predictions[0].NumDetections, res.Predictions[0].DetectionClasses)
	// fmt.Println("")

	predict := "NA"
	score := float32(0.0)
	categoryInd := int(res.Predictions[0].DetectionClasses[0])

	if res.Predictions[0].NumDetections > 0 {
		predict = cocoCategories[categoryInd]
		score = res.Predictions[0].DetectionScores[0]
	}

	modelRes := ModelResponse{
		Predict: predict,
		Score:   score,
	}

	return &modelRes, err
}

func parseImagenetResponse(response []byte) (*ModelResponse, error) {

	res := imagenetPredictResponseStruct{}

	err := json.Unmarshal(response, &res)

	// fmt.Printf("Unmarshalled Response: %+v\n ", res)

	// fmt.Printf("Response NumDetections: %f  Classes: %v ", res.Predictions[0].NumDetections, res.Predictions[0].DetectionClasses)
	// fmt.Println("")

	predict := res.ClassName
	score := res.Confidence

	modelRes := ModelResponse{
		Predict: predict,
		Score:   score,
	}

	return &modelRes, err
}

func parsePatternRecognitionResponse(response []byte) (*ModelResponse, error) {

	res := patternRecognitionPredictResponseStruct{}

	err := json.Unmarshal(response, &res)

	fmt.Printf("Unmarshalled Pattern RecognitionResponse: %+v\n ", res)

	predict := res.Prediction.Pattern
	score := res.Prediction.Max

	modelRes := ModelResponse{
		Predict: predict,
		Score:   score,
	}

	return &modelRes, err
}

func parseAnomalyDetectionResponse(response []byte) (*ModelResponse, error) {

	res := anomalyDetectionResponseStruct{}

	err := json.Unmarshal(response, &res)

	// predict := strconv.FormatFloat(res.ExpectedValue, 'f', 2, 64) + "|" + strconv.FormatBool(res.IsAnomaly)
	predict := strconv.FormatFloat(res.ExpectedValue, 'f', 2, 64)
	modelRes := ModelResponse{
		Predict: predict,
		Score:   0,
	}

	return &modelRes, err
}

// PredictImage - uses a model to predict an image
func predictImage(mi modelInfo, binaryValue []byte) (string, error) {

	fmt.Printf("Model URL: %+v\n ", mi.URL)

	encjson, _ := buildRequestFromBinary(mi, binaryValue)

	_, response := httpSender.Predict(loggingClient, mi.URL, string(encjson))

	bresponse, _ := response.([]byte)

	parsedResponse, _ := parseResponse(mi, bresponse)

	return parsedResponse.Predict, nil

}

// PredictPattern - uses a model to predict a pattern
func predictPattern(mi modelInfo, pattern string) (string, error) {

	fmt.Printf("Model URL: %+v\n ", mi.URL)

	encjson, _ := buildRequestFromPattern(mi, pattern)

	_, response := httpSender.Predict(loggingClient, mi.URL, string(encjson))

	bresponse, _ := response.([]byte)

	parsedResponse, _ := parseResponse(mi, bresponse)

	return parsedResponse.Predict, nil

}

// detectAnomaly- uses a model to detect an anomaly
func detectAnomaly(mi modelInfo, value float64) (string, error) {

	encjson, _ := buildAnomalyDetectionRequest(value)

	_, response := httpSender.Predict(loggingClient, mi.URL, string(encjson))

	bresponse, _ := response.([]byte)

	parsedResponse, _ := parseResponse(mi, bresponse)

	return parsedResponse.Predict, nil

}

// Predict - predicts the result for given reading
func Predict(reading models.Reading) (string, error) {

	// Get model information for device's resource
	mi := getModelInfo(reading.Device, reading.Name)

	prediction := ""

	switch {
	case reading.ValueType == "Binary":
		fmt.Println("Predicting Binary data")

		prediction, _ = predictImage(mi, reading.BinaryValue)

	case mi.ModelType == "pattern_recognition":
		fmt.Println("Predicting patterns")

		prediction, _ = predictPattern(mi, reading.Value)

	case mi.ModelType == "anomaly_detection":
		fmt.Println("Predicting anomalies")

		value, _ := strconv.ParseFloat(reading.Value, 64)

		prediction, _ = detectAnomaly(mi, value)
	default:
		fmt.Println("Predict - No model")
	}

	return prediction, nil
}
