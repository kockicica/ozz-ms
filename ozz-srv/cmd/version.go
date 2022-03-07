/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"runtime"

	"github.com/spf13/cobra"
)

var (
	NAME       = "gospl"
	VERSION    = "development version"
	REVISION   = "HEAD"
	BRANCH     = "HEAD"
	BUILT      = "unknown"
	AppVersion AppVersionInfo
)

type AppVersionInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Revision     string `json:"revision"`
	Branch       string `json:"branch"`
	GOVersion    string `json:"go_version"`
	BuiltAt      string `json:"built_at"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
}

func (v *AppVersionInfo) Line() string {
	return fmt.Sprintf("%s %s (%s)", v.Name, v.Version, v.Revision)
}

func (v *AppVersionInfo) ShortLine() string {
	return fmt.Sprintf("%s (%s)", v.Version, v.Revision)
}

func (v *AppVersionInfo) UserAgent() string {
	return fmt.Sprintf("%s %s (%s; %s; %s/%s)", v.Name, v.Version, v.Branch, v.GOVersion, v.OS, v.Architecture)
}

func (v *AppVersionInfo) Extended() string {
	version := fmt.Sprintf("Version:      %s\n", v.Version)
	version += fmt.Sprintf("Git revision: %s\n", v.Revision)
	version += fmt.Sprintf("Git branch:   %s\n", v.Branch)
	version += fmt.Sprintf("GO version:   %s\n", v.GOVersion)
	version += fmt.Sprintf("Built:        %s\n", v.BuiltAt)
	version += fmt.Sprintf("OS/Arch:      %s/%s\n", v.OS, v.Architecture)

	return version
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Get version",
	RunE: func(cmd *cobra.Command, args []string) error {
		short, err := cmd.Flags().GetBool("short")
		if err != nil {
			return err
		}
		if short {
			fmt.Println(AppVersion.ShortLine())
		} else {
			fmt.Println(AppVersion.Extended())
		}
		return nil
	},
}

func init() {

	AppVersion = AppVersionInfo{
		Name:         NAME,
		Version:      VERSION,
		Revision:     REVISION,
		Branch:       BRANCH,
		GOVersion:    runtime.Version(),
		BuiltAt:      BUILT,
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}

	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolP("short", "s", false, "Print short version information")

}
