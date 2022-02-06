package server

import (
	"errors"
	"net/http"

	"ozz-ms/pkg/data/model"
	"ozz-ms/pkg/data/repository"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Server) searchSchedules(ctx echo.Context) error {
	var err error
	ssp := repository.ScheduleSearchParams{}

	err = ctx.Bind(&ssp)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = ctx.Validate(ssp); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var data []model.Schedule
	if err := s.repo.Schedules(ssp, &data); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ret := []model.ScheduleDTO{}

	for _, sch := range data {
		ret = append(ret, sch.Map())
	}

	return ctx.JSON(http.StatusOK, ret)

}

func (s *Server) getSchedule(ctx echo.Context) error {

	var (
		id  int
		err error
	)

	err = echo.PathParamsBinder(ctx).Int("id", &id).BindError()

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sch := model.Schedule{}
	if err := s.repo.Schedule(id, &sch); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return ctx.JSON(http.StatusOK, sch.Map())

}

func (s *Server) createSchedule(ctx echo.Context) error {

	var err error

	data := model.NewScheduleDTO{}
	if err = ctx.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err = ctx.Validate(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	sch, err := s.repo.NewSchedule(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, sch.Map())
}

func (s *Server) deleteSchedule(ctx echo.Context) error {

	var (
		id  int
		err error
	)

	err = echo.PathParamsBinder(ctx).Int("id", &id).BindError()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = s.repo.DeleteSchedule([]int{id}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) updateSchedule(ctx echo.Context) error {

	var id int
	if err := echo.PathParamsBinder(ctx).Int("id", &id).BindError(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dto := model.NewScheduleDTO{}
	if err := ctx.Bind(&dto); err != nil {
		return err
	}

	if err := s.repo.SetSchedule(id, dto); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	sch := model.Schedule{}
	if err := s.repo.Schedule(id, &sch); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, sch.Map())
}
