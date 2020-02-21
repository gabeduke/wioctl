package sensors

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func getNewSensorGauge(name string) prometheus.Gauge {

	gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: fmt.Sprintf("wioctl_%s_gauge_current", name),
		Help: fmt.Sprintf("The current %s gauge reading", name),
	})

	return gauge

}
