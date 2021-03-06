package server

import (
	"net/http"
	"time"

	"ozz-ms/pkg/data/model"

	"github.com/labstack/echo/v4"
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
	return ctx.JSON(http.StatusOK, fnd)
}

//func (s *Server) increaseDispositionPlayedCount(ctx echo.Context) error {
//
//	var id int
//	if err := echo.PathParamsBinder(ctx).Int("id", &id).BindError(); err != nil {
//		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
//	}
//
//	time := time.Now()
//
//	dsp, err := s.repo.ChangePlayCountForDisposition(id, time, 1)
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return echo.NewHTTPError(http.StatusNotFound, err.Error())
//		}
//		return err
//	}
//
//	return ctx.JSON(http.StatusOK, dsp.Map())
//}
//
//func (s *Server) decreaseDispositionPlayedCount(ctx echo.Context) error {
//
//	var id int
//	if err := echo.PathParamsBinder(ctx).Int("id", &id).BindError(); err != nil {
//		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
//	}
//
//	time := time.Now()
//
//	dsp, err := s.repo.ChangePlayCountForDisposition(id, time, -1)
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return echo.NewHTTPError(http.StatusNotFound, err.Error())
//		}
//		return err
//	}
//
//	return ctx.JSON(http.StatusOK, dsp.Map())
//}

func (s *Server) createDispositions(ctx echo.Context) error {

	var cdp model.CreateDispositionParams
	if err := ctx.Bind(&cdp); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	data := []model.Schedule{}
	if err := s.repo.CreateDispositions(cdp, &data); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	res := []model.ScheduleDTO{}
	for _, schedule := range data {
		res = append(res, schedule.Map())
	}

	return ctx.JSON(http.StatusCreated, res)
}

func (s *Server) markDispositionExecution(ctx echo.Context) error {

	var ep model.DispositionExecuteParams
	if err := ctx.Bind(&ep); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := s.repo.MarkDispositionExecute(ep); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
