package wioctl

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Schedule int `mapstructure:"schedule,omitempty"`
	Wio      struct {
		BasePath string `mapstructure:"base_path"`
		BasePort string `mapstructure:"base_port"`
	} `mapstructure:"wio"`
	Sensors []Sensor `mapstructure:"sensors"`
}

type Sensor struct {
	Name            string `mapstructure:"name"`
	SensorTarget    string `mapstructure:"sensor_target"`
	SensorPath      string `mapstructure:"sensor_path"`
	SensorKey       string `mapstructure:"sensor_key"`
	SensorTokenKey  string `mapstructure:"sensor_token_key"`
	SensorTokenPath string `mapstructure:"sensor_token_path"`
}

// Config returns a safe copy of CivoCtl cfg
func (w *Wioctl) Config() *Config {
	w.lock.Lock()
	cfg := w.cfg
	w.lock.Unlock()
	return cfg
}

// SetConfig sets a safe copy of CivoCtl cfg
func (w *Wioctl) SetConfig(cfg *Config) {
	w.lock.Lock()
	w.cfg = cfg
	w.lock.Unlock()
}

// LoadConfig returns a config from viper and updates channel
func LoadConfig() (*Config, chan *Config) {

	config := &Config{}
	viper.Unmarshal(config)

	viper.WatchConfig()

	configCh := make(chan *Config, 1)

	prev := time.Now()
	viper.OnConfigChange(func(e fsnotify.Event) {
		now := time.Now()
		// fsnotify sometimes fires twice
		if now.Sub(prev) > time.Second {
			config := &Config{}
			err := viper.Unmarshal(config)
			if err == nil {
				configCh <- config
			}

			prev = now
		}
	})

	return config, configCh
}
