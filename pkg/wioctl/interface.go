package wioctl

import (
	"github.com/gabeduke/wioctl/pkg/sensors"
	"github.com/influxdata/influxdb-client-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

// NewWioctl generates a populated wioctler
func NewWioctl(cfg *Config, influxAddr string, influxToken string) (*Wioctl, error) {

	// create wioctl
	w := Wioctl{}

	// hang config
	w.cfg = cfg

	// hang influx client
	client := http.Client{}
	influx, err := influxdb.New(influxAddr, influxToken, influxdb.WithHTTPClient(&client))
	if err != nil {
		return nil, err
	}

	w.influxClient = influx

	return &w, nil
}

// Wioctler is the interface for wioctl
type Wioctler interface {
	Run()
}

// Wioctl is the definition for a wioctler
type Wioctl struct {
	cfg          *Config
	influxClient *influxdb.Client
	lock         sync.Mutex
}

func (w *Wioctl) Run() {

	log.WithFields(log.Fields{
		"path":     w.cfg.Wio.BasePath,
		"port":     w.cfg.Wio.BasePort,
		"schedule": w.cfg.Schedule,
	}).Debug("wio base params")

	sensors := w.cfg.Sensors
	for _, sensor := range sensors {
		err := w.generateHandler(w.cfg.Wio.BasePath, w.cfg.Wio.BasePort, sensor)
		if err != nil {
			log.Error(err)
		}
	}

	select {}
}

func (w *Wioctl) generateHandler(baseUrl string, port string, sensor Sensor) error {

	log.Infof("Generate Handler: %s", sensor.Name)

	s := &sensors.Wio{}

	switch sensor.Name {
	case "fahrenheit":
		s = sensors.NewFahrenheit(baseUrl, port, sensor.SensorTokenKey)
	case "celsius":
		s = sensors.NewCelsius(baseUrl, port, sensor.SensorTokenKey)
	case "lux":
		s = sensors.NewLux(baseUrl, port, sensor.SensorTokenKey)
	case "humidity":
		s = sensors.NewHumidity(baseUrl, port, sensor.SensorTokenKey)
	case "moisture":
		s = sensors.NewMoisture(baseUrl, port, sensor.SensorTokenKey)
	case "quality":
		s = sensors.NewAirQuality(baseUrl, port, sensor.SensorTokenKey)
	}

	tick := time.NewTicker(time.Second * time.Duration(w.cfg.Schedule))
	go w.scheduler(s, tick)

	return nil
}

func (w *Wioctl) scheduler(s *sensors.Wio, tick *time.Ticker) {
	for _ = range tick.C {
		reading, err := s.GetReading()
		if err != nil {
			log.Error(err)
		}

		s.ReportMetricInflux(w.influxClient, reading)

	}
}
