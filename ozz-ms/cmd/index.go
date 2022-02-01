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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ozz-ms/pkg/media_index"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	OVERWRITE_FLAG_NAME = "overwrite"
)

// indexCmd represents the media_index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Search index management",
	Long:  `Create or query media index`,
}

var createCmd = &cobra.Command{
	Use:   "create <folder-to-media>",
	Short: "Create search index",
	Long:  `Use command to create audio media search index, specifying locations to include in search`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		indexName := viper.GetString(INDEX_NAME_FLAG_NAME)
		verbose := viper.GetBool(VERBOSE_FLAG_NAME)

		if indexName == "" {
			return errors.New("no media_index name specified")
		}

		absIndexPath, err := filepath.Abs(indexName)
		if err != nil {
			return err
		}
		cmd.Println("Create index:", absIndexPath)

		overwrite, err := cmd.Flags().GetBool(OVERWRITE_FLAG_NAME)
		if err != nil {
			return err
		}
		if overwrite {
			// remove existing media_index, if any
			if _, err := os.Stat(absIndexPath); !os.IsNotExist(err) {
				if err = os.RemoveAll(absIndexPath); err != nil {
					return err
				}
			}
		}

		var progressBar *progressbar.ProgressBar
		if !verbose {
			progressBar = progressbar.NewOptions(-1,
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionSetDescription("[cyan]Indexing...[reset]"),
				progressbar.OptionShowCount(),
				progressbar.OptionShowIts(),
				progressbar.OptionSetItsString("items"),
				progressbar.OptionOnCompletion(func() {
					fmt.Printf("\n")
				}),
				progressbar.OptionSpinnerType(14),
				progressbar.OptionFullWidth(),
				progressbar.OptionThrottle(65*time.Millisecond),
				//progressbar.OptionClearOnFinish(),
			)
		}
		index := media_index.NewIndex(absIndexPath)
		err = index.Create()
		if err != nil {
			return err
		}

		audioFiles, err := media_index.NewAudioWalker(args)
		if err != nil {
			return err
		}

		//go func() {
		//	for {
		//		select {
		//		case count := <-index.BatchWritten:
		//			fmt.Printf("%d items indexed\n", count)
		//		}
		//	}
		//}()

		doit := true
		numberOfDocs := 0
		for doit {
			select {
			case audioFile, more := <-audioFiles.File:
				if !more {
					doit = false
					break
				}
				printVerbose(cmd, "Added:", audioFile.Path)
				if err = index.AddItem(audioFile); err != nil {
					doit = false
					break
				}
				if progressBar != nil {
					_ = progressBar.Add(1)
				}
				numberOfDocs++
			case <-audioFiles.Finished:
				doit = false
				break
			case <-audioFiles.Progress:
				break
			}
		}

		if progressBar != nil {
			_ = progressBar.Close()
		}

		err = index.Flush()
		if err != nil {
			return err
		}
		err = index.Close()
		if err != nil {
			return err
		}
		cmd.Println("MediaIndex created, documents:", numberOfDocs)
		return nil
	},
}

func printVerbose(cmd *cobra.Command, message ...interface{}) {
	verbose, err := cmd.Flags().GetBool(VERBOSE_FLAG_NAME)
	if err == nil {
		if verbose {
			cmd.Println(message...)
		}
	}
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Args:  cobra.MinimumNArgs(1),
	Short: "Query media search index",
	Long:  `Query for media index`,
	RunE: func(cmd *cobra.Command, args []string) error {

		indexName := viper.GetString(INDEX_NAME_FLAG_NAME)
		verbose := viper.GetBool(VERBOSE_FLAG_NAME)

		index := media_index.NewIndex(indexName)
		index.Verbose = verbose
		if err := index.Open(); err != nil {
			return err
		}
		queryString := strings.Join(args, " ")
		res, err := index.Query(queryString)
		if err != nil {
			return err
		}
		res.WriteOut()
		cmd.Println("Total:", len(res), " found.")
		if err = index.Close(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	indexCmd.AddCommand(createCmd)
	indexCmd.AddCommand(queryCmd)

	createCmd.Flags().Bool(OVERWRITE_FLAG_NAME, false, "overwrite media index if exists")
	rootCmd.AddCommand(indexCmd)

	viper.BindPFlags(indexCmd.PersistentFlags())

}
