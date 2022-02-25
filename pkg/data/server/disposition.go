package server

import (
	"errors"
	"net/http"
	"time"

	"ozz-ms/pkg/data/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Server) searchDispositions(ctx echo.Context) error {

	var err error

	sp := model.DispositionSearchParams{}

	if err = ctx.Bind(&sp); err != nil {
		return err
	}

	if err = ctx.Validate(&sp); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	parsedDate, err := time.Parse("2006-01-02", sp.Date)
	if err != nil {
		return err
	}
	fnd, err := s.repo.DispositionForShiftAndData(sp.Shift, parsedDate)
	if err != nil {
		return err
	}
	res := []model.DispositionDTO{}
	for _, dsp := range fnd {
		res = append(res, dsp.Map())
	}

	return ctx.JSON(http.StatusOK, res)
}

func (s *Server) increaseDispositionPlayedCount(ctx echo.Context) error {

	var id int
	if err := echo.PathParamsBinder(ctx).Int("id", &id).BindError(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	time := time.Now()

	dsp, err := s.repo.ChangePlayCountForDisposition(id, time, 1)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return ctx.JSON(http.StatusOK, dsp.Map())
}

func (s *Server) decreaseDispositionPlayedCount(ctx echo.Context) error {

	var id int
	if err := echo.PathParamsBinder(ctx).Int("id", &id).BindError(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	time := time.Now()

	dsp, err := s.repo.ChangePlayCountForDisposition(id, time, -1)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return ctx.JSON(http.StatusOK, dsp.Map())
}

func (s *Server) createDispositions(ctx echo.Context) error {

	var activeSchedules []model.Schedule

	if err := s.repo.ActiveSchedules(&activeSchedules); err != nil {
		return err
	}

	res := []model.ScheduleDTO{}
	for _, activeSchedule := range activeSchedules {
		if err := s.repo.CreateDispositions(&activeSchedule); err != nil {
			return err
		}
		res = append(res, activeSchedule.Map())
	}

	return ctx.JSON(http.StatusCreated, res)
}
