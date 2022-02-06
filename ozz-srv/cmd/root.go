/*
Copyright © 2022 kockicica@gmail.com

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
	"path/filepath"
	"syscall"

	"ozz-ms/pkg/data/server"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	serviceName        = "OZZZZZZMS"
	serviceDisplayName = "OZZZZZZ Media Server"
	serviceDescription = "OZZ replacement media media_index / file server"
)

var createdService service.Service
var runner *serverWrapper

type serverWrapper struct {
	service service.Service
	server  *server.Server
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
	absContentRoot, err := filepath.Abs(w.server.Config.RootPath)
	w.logger.Infof("Starting service, port: %d, dsn: %s, content root:%s", w.server.Config.Port, w.server.Config.Dsn, absContentRoot)
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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ozz-srv",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {

		cfg := server.ServerConfig{
			Dsn: "data.db",
			//Dsn:      "mysql://root:pass@localhost:3306/ozz?charset=utf8mb4&parseTime=True&loc=Local",
			Port:     27000,
			RootPath: "./media",
			Verbose:  true,
		}

		srv, err := server.NewDataServer(cfg)

		if err != nil {
			return err
		}

		runner := &serverWrapper{
			server: srv,
		}
		createdService, err = service.New(runner, &service.Config{
			Name:        serviceName,
			DisplayName: serviceDisplayName,
			Description: serviceDescription,
			Arguments:   []string{},
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ozz-srv.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ozz-srv" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ozz-srv")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
