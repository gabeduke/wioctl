package sensors

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var (
	fahrenheitGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wioctl_fahrenheit_gauge_current",
		Help: "The current fahrenheit gauge reading",
	})
)

type fahrenheitJson struct {
	FahrenheitDegree json.Number `json:"fahrenheit_degree"`
}

func NewFahrenheit(baseURL string, port string, tokenKey string) *Wio {

	s := &Sensor{
		id:      "GroveTempHumD0",
		path:    "temperature_f",
		handler: fahrenheitHandler,
	}

	w := NewSensor("fahrenheit", s, baseURL, port, tokenKey)

	return w
}

func fahrenheitHandler(logger *log.Entry, response *http.Response) float64 {
	body, _ := ioutil.ReadAll(response.Body)
	logger.Info(string([]byte(body)))

	defer response.Body.Close()

	var f fahrenheitJson
	err := json.Unmarshal([]byte(body), &f)
	if err != nil {
		logger.Error(err)
	}

	value, err := f.FahrenheitDegree.Float64()
	if err != nil {
		return 0
	}

	fahrenheitGauge.Set(value)

	return value
}
