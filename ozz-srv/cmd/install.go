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
	"os"
	"path/filepath"

	"ozz-ms/pkg/data/server"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   "Install ozzzz-srv as system service",
	RunE: func(cmd *cobra.Command, args []string) error {

		cmd.Println("Installing service")

		cfg := createServerConfig()
		absRoot, err := filepath.Abs(cfg.RootPath)
		if err != nil {
			return err
		}
		cfg.RootPath = absRoot

		// get current working dir
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		if cfg.Dsn == "data.db" {
			// without proper database url we need abs path
			cfg.Dsn = filepath.Join(wd, "data.db")
			cfg.Dsn = fmt.Sprintf("sqlite:///%s", cfg.Dsn)
		}
		srv, err := server.NewDataServer(cfg)

		if err != nil {
			return err
		}

		runner := &serverWrapper{
			server: srv,
		}

		serviceCfg := defaultServiceConfig()

		if service.Platform() == "windows-service" {
			serviceCfg.UserName = "Nt Authority\\Network service"
		}

		serviceCfg.Arguments = []string{
			"--database", cfg.Dsn,
			"--root", cfg.RootPath,
			"--port", fmt.Sprintf("%d", cfg.Port),
		}

		createdService, err = service.New(runner, serviceCfg)
		if err != nil {
			return err
		}

		err = service.Control(createdService, "install")
		if err != nil {
			return err
		}
		cmd.Println("Service installed.")
		return nil
	},
}

var uninstallCmd = &cobra.Command{
	Use:     "uninstall",
	Aliases: []string{"u"},
	Short:   "Uninstall server as system service",
	RunE: func(cmd *cobra.Command, args []string) error {

		cmd.Println("Uninstalling service")

		cfg := createServerConfig()
		srv, err := server.NewDataServer(cfg)

		if err != nil {
			return err
		}

		runner := &serverWrapper{
			server: srv,
		}

		serviceCfg := defaultServiceConfig()

		createdService, err = service.New(runner, serviceCfg)
		if err != nil {
			return err
		}

		err = service.Control(createdService, "uninstall")
		if err != nil {
			return err
		}

		cmd.Println("Service uninstalled")

		return nil
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start server system service",
	RunE: func(cmd *cobra.Command, args []string) error {

		cmd.Println("Starting service")
		cfg := createServerConfig()
		srv, err := server.NewDataServer(cfg)

		if err != nil {
			return err
		}

		runner := &serverWrapper{
			server: srv,
		}

		serviceCfg := defaultServiceConfig()

		createdService, err = service.New(runner, serviceCfg)
		if err != nil {
			return err
		}

		err = service.Control(createdService, "start")
		if err != nil {
			return err
		}
		cmd.Println("Service started")
		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop server system service",
	RunE: func(cmd *cobra.Command, args []string) error {

		cmd.Println("Stopping service")
		cfg := createServerConfig()
		srv, err := server.NewDataServer(cfg)

		if err != nil {
			return err
		}

		runner := &serverWrapper{
			server: srv,
		}

		serviceCfg := defaultServiceConfig()

		createdService, err = service.New(runner, serviceCfg)
		if err != nil {
			return err
		}

		err = service.Control(createdService, "stop")
		if err != nil {
			return err
		}
		cmd.Println("Service stopped")
		return nil
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart server system service",
	RunE: func(cmd *cobra.Command, args []string) error {

		cmd.Println("Restarting service")
		cfg := createServerConfig()
		srv, err := server.NewDataServer(cfg)

		if err != nil {
			return err
		}

		runner := &serverWrapper{
			server: srv,
		}

		serviceCfg := defaultServiceConfig()

		createdService, err = service.New(runner, serviceCfg)
		if err != nil {
			return err
		}

		err = service.Control(createdService, "restart")
		if err != nil {
			return err
		}
		cmd.Println("Service restarted")
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get service status",
	RunE: func(cmd *cobra.Command, args []string) error {

		cfg := createServerConfig()
		srv, err := server.NewDataServer(cfg)

		if err != nil {
			return err
		}

		runner := &serverWrapper{
			server: srv,
		}

		serviceCfg := defaultServiceConfig()

		createdService, err = service.New(runner, serviceCfg)
		if err != nil {
			return err
		}
		status, err := createdService.Status()
		if err != nil {
			return err
		}
		switch status {
		case service.StatusRunning:
			fmt.Println("Service is running")
		case service.StatusStopped:
			fmt.Println("Service is stopped")
		case service.StatusUnknown:
			fmt.Println("Service status is unknown, or service is not installed")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(statusCmd)
}
