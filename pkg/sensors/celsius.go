package sensors

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type celsiusJson struct {
	CelsiusDegree json.Number `json:"celsius_degree"`
}

func NewCelsius(baseURL string, port string, tokenKey string) *Wio {

	s := &Sensor{
		id:      "GroveTempHumD0",
		path:    "temperature",
		handler: celsiusHandler,
	}

	w := NewSensor("celsius", s, baseURL, port, tokenKey)

	return w
}

func celsiusHandler(logger *log.Entry, response *http.Response) float64 {
	body, _ := ioutil.ReadAll(response.Body)
	logger.Info(string([]byte(body)))

	defer response.Body.Close()

	var f celsiusJson
	err := json.Unmarshal([]byte(body), &f)
	if err != nil {
		logger.Error(err)
	}

	value, err := f.CelsiusDegree.Float64()
	if err != nil {
		return 0
	}

	gauge := getNewSensorGauge("celsius")
	gauge.Set(value)


	return value
}
