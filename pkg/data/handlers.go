package data

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Server) getShifts(ctx echo.Context) error {

	var data []struct {
		ID     int
		Name   string
		Active bool
	}

	if err := s.repo.Shifts(&data); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(200, data)

}

func (s *Server) getCategories(ctx echo.Context) error {

	var cats []CategoryDTO

	if err := s.repo.Categories(&cats); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, cats)
}

func (s *Server) authorize(ctx echo.Context) error {

	var login = struct {
		Username string `validate:"required" json:"username"`
		Password string `json:"password"`
	}{}

	if err := ctx.Bind(&login); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(&login); err != nil {
		return err
	}

	var user struct {
		Username string `json:"username"`
		Level    int    `json:"level"`
	}
	err := s.repo.Authorize(login.Username, login.Password, &user)
	switch err {
	case gorm.ErrRecordNotFound:
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case nil:
		return ctx.JSON(http.StatusOK, user)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
