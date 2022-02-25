package server

import (
	"github.com/gookit/validate"
)

type ServerValidator struct {
}

func (sv *ServerValidator) Validate(i interface{}) error {

	v := validate.Struct(i)

	if !v.Validate() {
		//return echo.NewHTTPError(http.StatusBadRequest, v.Errors)
		return v.Errors
	}
	return nil
}

func NewServerValidator() *ServerValidator {
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false

		//opt.SkipOnEmpty = false
		//opt.CheckDefault = false

		//opt.CheckZero = true
	})
	v := new(ServerValidator)
	return v
}
