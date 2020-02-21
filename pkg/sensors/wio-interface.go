package sensors

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"time"
)

type WioInterface interface {
	GetReading() (interface{}, error)
}

type Wio struct {
	url    *url.URL
	req    *http.Request
	Sensor *Sensor
	logger *log.Entry
}

type Sensor struct {
	id      string
	path    string
	handler func(logger *log.Entry, response *http.Response) float64
}

func NewSensor(name string, sensor *Sensor, baseUrl string, port string, tokenKey string) *Wio {

	c := Wio{
		Sensor: sensor,
	}

	c.logger = log.WithFields(log.Fields{
		"sensor": name,
	})

	c.configureURL(baseUrl, port)
	c.configureRequest(tokenKey)

	return &c
}

func (w *Wio) configureURL(baseUrl string, port string) {
	w.logger.Debug("Configure URL")

	u := &url.URL{
		Host:   fmt.Sprintf("%s:%s", baseUrl, port),
		Scheme: "https",
		Path:   fmt.Sprintf("v1/node/%s/%s", w.Sensor.id, w.Sensor.path),
	}

	w.url = u
}

func (w *Wio) configureRequest(tokenKey string) {
	w.logger.Debug("Configure Request")

	// Create a Bearer string by appending string access token
	var bearer = "token " + os.Getenv(tokenKey)
	w.logger.WithFields(log.Fields{
		"URL":           w.url.String(),
		"Authorization": bearer,
	}).Trace("Request")

	// Create a new request using http
	req, err := http.NewRequest("GET", w.url.String(), nil)
	if err != nil {
		w.logger.Error(err)
	}

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	w.req = req
}

func (w *Wio) GetReading() (float64, error) {
	w.logger.Debug("Get Reading")

	client := &http.Client{}
	w.logger.Trace(w.req)

	resp, err := client.Do(w.req)
	if err != nil {
		return 0, nil
	}

	value := w.Sensor.handler(w.logger, resp)

	return value, nil

}

func (w *Wio) ReportMetricPromPushGateway(influxClient *influxdb.Client, value float64) error {
	w.logger.Debug("Report Metric")

	if value == 0 {
		w.logger.Errorf("no value found for %s", w.Sensor.id)
	}

	metrics := []influxdb.Metric{
		influxdb.NewRowMetric(
			map[string]interface{}{w.Sensor.id: value},
			"fleet-metrics",
			map[string]string{
				"sensor_path": w.Sensor.path,
			},
			time.Now().UTC()),
	}

	_, err := influxClient.Write(context.Background(), "Fleet IOT", "c670b60f97bc7205", metrics...)
	if err != nil {
		return err
	}

	w.logger.Infof("Influx write successful")
	return nil
}

func (w *Wio) ReportMetricInflux(influxClient *influxdb.Client, value float64) error {
	w.logger.Debug("Report Metric")

	if value == 0 {
		w.logger.Errorf("no value found for %s", w.Sensor.id)
	}

	metrics := []influxdb.Metric{
		influxdb.NewRowMetric(
			map[string]interface{}{w.Sensor.id: value},
			"fleet-metrics",
			map[string]string{
				"sensor_path": w.Sensor.path,
			},
			time.Now().UTC()),
	}

	_, err := influxClient.Write(context.Background(), "Fleet IOT", "c670b60f97bc7205", metrics...)
	if err != nil {
		return err
	}

	w.logger.Infof("Influx write successful")
	return nil
}

