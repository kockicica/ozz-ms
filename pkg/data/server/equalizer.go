package server

import (
	"errors"
	"net/http"

	"ozz-ms/pkg/data/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Server) getEqualizers(ctx echo.Context) error {

	var name string

	if err := echo.QueryParamsBinder(ctx).String("name", &name).BindError(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if name == "" {
		data := []model.Equalizer{}
		if err := s.repo.Equalizers(&data); err != nil {
			return err
		}

		res := []model.EqualizerDTO{}
		for _, equalizer := range data {
			res = append(res, equalizer.Map())
		}

		return ctx.JSON(http.StatusOK, res)
	} else {
		data := model.Equalizer{}
		if err := s.repo.EqualizerByName(name, &data); err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return ctx.JSON(http.StatusOK, data.Map())
	}

}

func (s *Server) createEqualizer(ctx echo.Context) error {

	eqd := model.EqualizerDTO{}
	if err := ctx.Bind(&eqd); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(&eqd); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	eq := model.Equalizer{
		Name:   eqd.Name,
		PreAmp: eqd.PreAmp,
		Amp1:   eqd.Amp1,
		Amp2:   eqd.Amp2,
		Amp3:   eqd.Amp3,
		Amp4:   eqd.Amp4,
		Amp5:   eqd.Amp5,
		Amp6:   eqd.Amp6,
		Amp7:   eqd.Amp7,
		Amp8:   eqd.Amp8,
		Amp9:   eqd.Amp9,
		Amp10:  eqd.Amp10,
	}

	if err := s.repo.NewEqualizer(&eq); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, eq.Map())
}

func (s *Server) updateEqualizer(ctx echo.Context) error {

	var id int
	if err := echo.PathParamsBinder(ctx).Int("id", &id).BindError(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dto := model.EqualizerDTO{}
	if err := ctx.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := s.repo.SetEqualizer(id, dto); err != nil {
		return err
	}
	dto.ID = uint(id)

	return ctx.JSON(http.StatusOK, dto)
}

func (s *Server) deleteEqualizer(ctx echo.Context) error {

	var id int

	if err := echo.PathParamsBinder(ctx).Int("id", &id).BindError(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := s.repo.DeleteEqualizer(id); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) getEqualizer(ctx echo.Context) error {

	var id int
	if err := echo.PathParamsBinder(ctx).Int("id", &id).BindError(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var d model.Equalizer
	if err := s.repo.Equalizer(id, &d); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return ctx.JSON(http.StatusOK, d.Map())

}
