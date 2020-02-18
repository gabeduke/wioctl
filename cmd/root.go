/*
Copyright © 2020 Gabriel Duke

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

const app = "wioctl"

var (
	cfgFile string
	VERSION string
)

var rootCmd = &cobra.Command{
	Use:   app,
	Short: "Wio logger controller",
	Long: `The wio control loop will watch a given list of devices
and log the sensor data to influxdb.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  `All software has versions. This is Wio's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION)
	},
}

func Execute(version string) {
	VERSION = version
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(versionCmd)

	flags := rootCmd.PersistentFlags()

	flags.StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/.%s.yaml)", app))
	flags.String("influx-token", "", "Influx API Token (env variable: WIOCTL_INFLUX_TOKEN)")
	flags.String("influx-addr", "", "Influx Address (env variable: WIOCTL_INFLUX_ADDR)")
	flags.StringP("log-level", "l", "info", "Log level (error|INFO|debug|trace)")
	flags.IntP("schedule", "s", 60, "Schedule (in seconds)")

	viper.BindPFlags(flags)

	viper.SetEnvPrefix("wioctl")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.BindEnv("influx-token")
	viper.BindEnv("influx-addr")
}

func initLogger() {
	level := viper.GetString("log-level")
	log.Infof("Log level: %s", level)

	switch level {
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
		viper.Debug()
	case "trace":
		log.SetLevel(log.TraceLevel)
		viper.Debug()
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Error(err)
		}

		// Search config in home directory with name ".civoctl" (without extension).
		viper.AddConfigPath(fmt.Sprintf("/etc/%s/", app))
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")

		viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
		viper.SetConfigName("." + app)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	} else {
		log.Fatal("config file not found")
	}

}
