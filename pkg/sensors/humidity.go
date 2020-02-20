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
	humidityGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wioctl_humidity_gauge_current",
		Help: "The current humidity gauge reading",
	})
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

	humidityGauge.Set(value)

	return value
}
