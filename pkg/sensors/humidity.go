package sensors

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type humidityJson struct {
	Humidity json.Number `json:"humidity"`
}

func NewHumidity(baseURL string, port string, tokenKey string) *Wio {

	s := &Sensor{
		id:      "GroveTempHumD0",
		path:    "humidity",
		handler: humidityHandler,
	}

	w := NewSensor("humidity", s, baseURL, port, tokenKey)

	return w
}

func humidityHandler(logger *log.Entry, response *http.Response) float64 {
	body, _ := ioutil.ReadAll(response.Body)
	logger.Info(string([]byte(body)))

	defer response.Body.Close()

	var f humidityJson
	err := json.Unmarshal([]byte(body), &f)
	if err != nil {
		logger.Error(err)
	}

	value, err := f.Humidity.Float64()
	if err != nil {
		return 0
	}

	gauge := getNewSensorGauge("humidity")
	gauge.Set(value)


	return value
}
