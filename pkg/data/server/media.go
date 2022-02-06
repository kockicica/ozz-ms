package server

import (
	"errors"
	"net/http"

	"ozz-ms/pkg/data/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Server) serveAudioFile(ctx echo.Context) error {

	var id int
	if err := echo.PathParamsBinder(ctx).Int("id", &id).BindError(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	rec := model.AudioRecording{}
	if err := s.repo.AudioRecording(id, &rec); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	//absPath, err := filepath.Abs(rec.Path)
	//if err != nil {
	//	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	//}

	//f, err := os.Open(rec.Path)
	//if err != nil {
	//	return err
	//}
	//return ctx.Stream(200, "audio/mpeg", f)

	return ctx.File(rec.Path)

}
