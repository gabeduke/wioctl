package sensors

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type moistureJson struct {
	Moisture json.Number `json:"analog"`
}

func NewMoisture(baseURL string, port string, tokenKey string) *Wio {

	s := &Sensor{
		id:      "GenericAInA0",
		path:    "analog",
		handler: moistureHandler,
	}

	w := NewSensor("moisture", s, baseURL, port, tokenKey)

	return w
}

func moistureHandler(logger *log.Entry, response *http.Response) float64 {
	body, _ := ioutil.ReadAll(response.Body)
	logger.Info(string([]byte(body)))

	defer response.Body.Close()

	var f moistureJson
	err := json.Unmarshal([]byte(body), &f)
	if err != nil {
		logger.Error(err)
	}

	value, err := f.Moisture.Float64()
	if err != nil {
		return 0
	}

	return value
}
