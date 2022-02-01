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
	"os"

	"ozz-ms/pkg/media_index"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	VERBOSE_FLAG_NAME    = "verbose"
	INDEX_NAME_FLAG_NAME = "index-name"
	PORT_FLAG_NAME       = "port"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "ozz-ms",
	Short: "OZZ media server command line interface",
	Long:  `OZZ media server command line interface`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ozz-ms.yaml)")
	rootCmd.PersistentFlags().BoolP(VERBOSE_FLAG_NAME, "v", false, "display verbose output")
	rootCmd.PersistentFlags().String(INDEX_NAME_FLAG_NAME, media_index.DefaultIndexName, "index name")
	viper.BindPFlags(rootCmd.PersistentFlags())
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

		// Search config in home directory with name ".ozz-ms-ozz-ms" (without extension).
		viper.AddConfigPath(home)
		viper.SetEnvPrefix("OZZ_MS_")
		viper.SetConfigName(".ozz-ms")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
