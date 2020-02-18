/*
Copyright © 2020 Gabriel Duke <gabeduke@gmail.com>

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
	"github.com/gabeduke/wioctl/pkg/wioctl"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the wioctl control loop",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLogger()
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Beginning Wioctl control loop")

		token := viper.GetString("influx-token")
		addr := viper.GetString("influx-addr")

		c, cfgCh := wioctl.LoadConfig()
		app, err := wioctl.NewWioctl(c, addr, token)
		if err != nil {
			log.Error(err)
		}

		go func() {
			for {
				app.SetConfig(<-cfgCh)
				log.Info("NewWioctl config loaded")
			}
		}()

		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
