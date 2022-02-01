package server

import (
	"context"
	"fmt"
	"time"

	"ozz-ms/pkg/media_index"

	"github.com/labstack/echo/v4"
)

type OzzServerConfig struct {
	IndexName string
	Port      int
	Verbose   bool
}

type OzzServer struct {
	Config OzzServerConfig
	e      *echo.Echo
	index  *media_index.MediaIndex
}

func (s *OzzServer) Start() error {

	err := s.index.Open()
	if err != nil {
		return err
	}
	err = s.e.Start(fmt.Sprintf(":%d", s.Config.Port))
	if err != nil {
		return err
	}

	return nil
}

func (s *OzzServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.e.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewOzzServer(config OzzServerConfig) *OzzServer {
	ozs := OzzServer{
		Config: config,
	}
	ozs.e = echo.New()
	ozs.e.HideBanner = true
	ozs.e.HidePort = true
	ozs.index = media_index.NewIndex(config.IndexName)
	ozs.e.GET("/media", ozs.searchMedia)
	ozs.e.GET("/media/:id", ozs.getMedia)
	ozs.e.GET("/media/stream/:id", ozs.getMediaStream)
	return &ozs
}
