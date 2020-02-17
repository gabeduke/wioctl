package wioctl

import (
	"github.com/influxdata/influxdb-client-go"
	"net/http"
	"sync"
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
	cfg *Config
	influxClient *influxdb.Client
	lock sync.Mutex
}

func (w *Wioctl) Run()  {

}