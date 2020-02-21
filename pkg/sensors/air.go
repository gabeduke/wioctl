package sensors

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type airQualityJson struct {
	Quality json.Number `json:"quality"`
}

func NewAirQuality(baseURL string, port string, tokenKey string) *Wio {

	s := &Sensor{
		id:      "GroveAirqualityA0",
		path:    "quality",
		handler: airQualityHandler,
	}

	w := NewSensor("air_quality", s, baseURL, port, tokenKey)

	return w
}

func airQualityHandler(logger *log.Entry, response *http.Response) float64 {
	body, _ := ioutil.ReadAll(response.Body)
	logger.Info(string([]byte(body)))

	defer response.Body.Close()

	var f airQualityJson
	err := json.Unmarshal([]byte(body), &f)
	if err != nil {
		logger.Error(err)
	}

	value, err := f.Quality.Float64()
	if err != nil {
		return 0
	}

	gauge := getNewSensorGauge("air")
	gauge.Set(value)


	return value
}
