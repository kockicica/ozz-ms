/*
Copyright © 2022 kockicica

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"ozz-ms/pkg/server"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	serviceName        = "OZZZZZZMS"
	serviceDisplayName = "OZZZZZZ Media Server"
	serviceDescription = "OZZ replacement media media_index / file server"
)

var createdService service.Service
var runner *serverWrapper

type serverWrapper struct {
	service service.Service
	server  *server.OzzServer
	exit    chan os.Signal
	logger  service.Logger
}

func (w *serverWrapper) Start(s service.Service) error {
	w.exit = make(chan os.Signal)
	w.service = s
	var err error
	w.logger, err = s.Logger(nil)
	if err != nil {
		return err
	}
	w.logger.Infof("Starting service, port: %d, media_index: %s", w.server.Config.Port, w.server.Config.IndexName)
	// starting service
	go w.run()
	// service started
	w.logger.Infof("Service started")
	return nil
}

func (w *serverWrapper) Stop(s service.Service) error {
	w.logger.Infof("Stopping service")
	err := w.server.Stop()
	if err != nil {
		return err
	}
	close(w.exit)
	// service stopped
	w.logger.Infof("Service stopped")
	return nil
}

func (w *serverWrapper) run() {
	signal.Notify(w.exit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		err := w.server.Start()
		if err != http.ErrServerClosed {
			// should log some error
			w.logger.Errorf("Error starting server:%s", err)
		} else {
			w.logger.Info("Server stopped")
		}
	}()

	<-w.exit
}

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Media query server management",
	Long:  `Use subcommands to manage media server (start / stop / install / uninstall)`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		port := viper.GetInt(PORT_FLAG_NAME)
		indexName := viper.GetString(INDEX_NAME_FLAG_NAME)

		cfg := server.OzzServerConfig{
			Port:      port,
			IndexName: indexName,
		}
		srv := server.NewOzzServer(cfg)

		runner = &serverWrapper{}
		runner.server = srv

		return nil
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run server",
	Long:  `Run audio media server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		createdService, err = service.New(runner, &service.Config{
			Name:        serviceName,
			DisplayName: serviceDisplayName,
			Description: serviceDescription,
			Arguments: []string{
				INDEX_NAME_FLAG_NAME, runner.server.Config.IndexName,
				"service", "run",
				PORT_FLAG_NAME, fmt.Sprintf("%d", runner.server.Config.Port),
			},
		})
		if err != nil {
			return err
		}

		if err := createdService.Run(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	serviceCmd.AddCommand(runCmd)

	rootCmd.AddCommand(serviceCmd)

	serviceCmd.PersistentFlags().IntP(PORT_FLAG_NAME, "p", 26000, "port to serve on")
	viper.BindPFlags(serviceCmd.PersistentFlags())

}
