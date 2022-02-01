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
	"log"
	"os"
	"runtime"
	"time"

	"ozz-ms/pkg/oto"

	"github.com/spf13/cobra"
	"github.com/tosone/minimp3"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play audio file",
	RunE: func(cmd *cobra.Command, args []string) error {

		f, err := os.Open("Z:\\documents\\My Music\\cistiliste\\Maya Jane Coles\\2021 - Maya Jane Coles - Night Creature\\03 - Maya Jane Coles - N31.mp3")
		if err != nil {
			return err
		}

		var dec *minimp3.Decoder
		dec, err = minimp3.NewDecoder(f)
		if err != nil {
			return err
		}

		started := dec.Started()
		<-started

		log.Printf("Convert audio sample rate: %d, channels: %d\n", dec.SampleRate, dec.Channels)

		var ctx *oto.Context
		var r chan struct{}
		ctx, r, err = oto.NewContextWithDevice(dec.SampleRate, dec.Channels, 2, 0)
		if err != nil {
			return err
		}
		<-r

		p := ctx.NewPlayer(dec)
		p.Play()
		p.SetVolume(0.2)
		up := p.UnplayedBufferSize()
		for up > 0 {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("Unplayed buffer size:", up)
			up = p.UnplayedBufferSize()
		}
		p.Close()
		runtime.KeepAlive(p)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
}
