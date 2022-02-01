package server

import (
	"errors"
	"os"

	"github.com/labstack/echo/v4"
)

func (s *OzzServer) getMedia(ctx echo.Context) error {
	id := ctx.Param("id")
	path := s.index.GetPath(id)
	if path == "" {
		return errors.New("unable to find media")
	}
	return ctx.File(path)
}

func (s *OzzServer) getMediaStream(ctx echo.Context) error {
	id := ctx.Param("id")
	path := s.index.GetPath(id)
	if path == "" {
		return errors.New("unable to find media")
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return ctx.Stream(200, "audio/mpeg", f)
}

func (s *OzzServer) searchMedia(ctx echo.Context) error {
	q := ctx.QueryParam("q")
	audioFiles, err := s.index.Query(q)
	if err != nil {
		return err
	}
	return ctx.JSON(200, audioFiles)
}
