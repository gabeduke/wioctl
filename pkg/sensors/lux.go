package sensors

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type luxJson struct {
	Lux json.Number `json:"lux,omitempty"`
	Err string      `json:"error,omitempty"`
}

func NewLux(baseURL string, port string, tokenKey string) *Wio {

	s := &Sensor{
		id:      "GroveDigitalLightI2C0",
		path:    "lux",
		handler: luxHandler,
	}

	w := NewSensor("lux", s, baseURL, port, tokenKey)

	return w
}

func luxHandler(logger *log.Entry, response *http.Response) float64 {
	body, _ := ioutil.ReadAll(response.Body)
	logger.Info(string([]byte(body)))

	defer response.Body.Close()

	var f luxJson
	err := json.Unmarshal([]byte(body), &f)
	if err != nil {
		logger.Error(err)
	}

	if f.Err != "" {
		logger.Error(fmt.Errorf(f.Err))
		return 0
	}

	value, err := f.Lux.Float64()
	if err != nil {
		return 0
	}

	gauge := getNewSensorGauge("lux")
	gauge.Set(value)

	return value
}
