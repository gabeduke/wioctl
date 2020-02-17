package wioctl

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

type Config struct {

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

	//TODO check for nil cfg

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

