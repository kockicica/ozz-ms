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
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	vlc "github.com/adrg/libvlc-go/v3"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AudioPlayer struct {
	Finished      chan os.Signal
	TimeChanged   chan time.Duration
	MediaPosition chan float32
	vlcPlayer     *vlc.Player
	media         *vlc.Media
	mgr           *vlc.EventManager
	eventIDs      []vlc.EventID
}

func NewAudioPlayer() (*AudioPlayer, error) {
	pl := new(AudioPlayer)
	pl.Finished = make(chan os.Signal)
	pl.TimeChanged = make(chan time.Duration)
	pl.MediaPosition = make(chan float32)

	err := vlc.Init("--no-video", "--quiet")
	if err != nil {
		return nil, err
	}

	pl.vlcPlayer, err = vlc.NewPlayer()
	if err != nil {
		return nil, err
	}

	pl.mgr, err = pl.vlcPlayer.EventManager()
	if err != nil {
		return nil, err
	}

	events := []vlc.Event{
		vlc.MediaPlayerTimeChanged,
		vlc.MediaPlayerEndReached,
	}
	pl.eventIDs = []vlc.EventID{}
	for _, event := range events {
		eventID, err := pl.mgr.Attach(event, pl.eventCallback, nil)
		if err != nil {
			return nil, err
		}
		pl.eventIDs = append(pl.eventIDs, eventID)
	}

	return pl, nil
}

func (p *AudioPlayer) eventCallback(event vlc.Event, userData interface{}) {
	switch event {
	case vlc.MediaPlayerTimeChanged:
		mediaTime, err := p.vlcPlayer.MediaTime()
		if err != nil {
			log.Println(err)
			break
		}
		duration, err := time.ParseDuration(fmt.Sprintf("%dms", mediaTime))
		if err != nil {
			log.Println(err)
			break
		}

		mediaPosition, err := p.vlcPlayer.MediaPosition()
		if err != nil {
			log.Println(err)
			break
		}

		select {
		case p.TimeChanged <- duration:
		case p.MediaPosition <- mediaPosition:
		default:

		}

	case vlc.MediaPlayerEndReached:
		p.Finished <- syscall.SIGTERM
	}

}

func (p *AudioPlayer) Release() error {
	var err error

	if p.media != nil {
		err = p.media.Release()
		if err != nil {
			return err
		}
	}
	//for _, eventID := range p.eventIDs {
	//	p.mgr.Detach(eventID)
	//}
	p.mgr.Detach(p.eventIDs...)

	if err = p.vlcPlayer.Stop(); err != nil {
		return err
	}
	if err = p.vlcPlayer.Release(); err != nil {
		return err
	}

	if err = vlc.Release(); err != nil {
		return err
	}

	close(p.MediaPosition)
	close(p.Finished)
	close(p.TimeChanged)

	return nil
}

func (p *AudioPlayer) PlayMediaFromPath(path string) error {
	var err error
	p.media, err = p.vlcPlayer.LoadMediaFromPath(path)
	if err != nil {
		return err
	}

	err = p.vlcPlayer.Play()
	if err != nil {
		return err
	}

	return nil
}

func (p *AudioPlayer) SetVolume(volume int) error {
	return p.vlcPlayer.SetVolume(volume)
}

func (p *AudioPlayer) Volume() (int, error) {
	return p.vlcPlayer.Volume()
}

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play audio file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		pl, err := NewAudioPlayer()
		if err != nil {
			return err
		}

		volume := viper.GetInt("volume")
		if err = pl.SetVolume(volume); err != nil {
			return err
		}

		err = pl.PlayMediaFromPath(args[0])
		if err != nil {
			return err
		}

		signal.Notify(pl.Finished, syscall.SIGINT, syscall.SIGTERM)

		_, fileName := filepath.Split(args[0])
		//pb := progressbar.Default(100, fileName)
		pb := progressbar.NewOptions(100,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetPredictTime(true),
			progressbar.OptionSetDescription(fmt.Sprintf("[yellow]%s[reset]", fileName)),
			progressbar.OptionFullWidth(),
			progressbar.OptionOnCompletion(func() {
				fmt.Println("")
			}),
			progressbar.OptionSetRenderBlankState(true),
		)

		done := false
		for !done {
			select {
			case <-pl.Finished:
				done = true
				break
			case _ = <-pl.TimeChanged:
				//cmd.Printf("Time: %v\n", d)
			case p := <-pl.MediaPosition:
				//cmd.Printf("Perc:%.0f\n", p*100)
				err = pb.Set(int(p * 100))
				if err != nil {
					return err
				}
			}
		}
		if err = pb.Close(); err != nil {
			return err
		}

		cmd.Println("play finished...")

		if err = pl.Release(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(playCmd)

	playCmd.Flags().IntP("volume", "v", 100, "set volume (0 - 100)")
	viper.BindPFlags(playCmd.Flags())

}
