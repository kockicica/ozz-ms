package server

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

type ServerBinding struct {
}

func (s *ServerBinding) Bind(i interface{}, c echo.Context) error {

	db := new(echo.DefaultBinder)
	if err := db.Bind(i, c); err != nil && err != echo.ErrUnsupportedMediaType {
		return err
	}

	ctype := c.Request().Header.Get(echo.HeaderContentType)
	if strings.HasPrefix(ctype, echo.MIMEApplicationForm) || strings.HasPrefix(ctype, echo.MIMEMultipartForm) {
		var form *multipart.Form
		form, err := c.MultipartForm()
		if err == nil {
			err = echoBindFile(i, c, form.File)
		}
		return err
	}

	return nil
}

func echoBindFile(i interface{}, ctx echo.Context, files map[string][]*multipart.FileHeader) error {
	ival := reflect.Indirect(reflect.ValueOf(i))
	if ival.Kind() != reflect.Struct {
		return fmt.Errorf("input is not a struct pointer, indirect type is %s", ival.Type().String())
	}

	itype := ival.Type()
	for i := 0; i < itype.NumField(); i++ {
		ftype := itype.Field(i)
		fvalue := ival.Field(i)
		if !fvalue.CanSet() {
			continue
		}
		switch ftype.Type {
		case typeMultipartFileHeader:
			file := getFiles(files, ftype.Name, ftype.Tag.Get("form"))
			if file != nil && len(file) > 0 {
				fvalue.Set(reflect.ValueOf(file[0]))
			}
		case typeMultipartSliceFileHeader:
			file := getFiles(files, ftype.Name, ftype.Tag.Get("form"))
			if file != nil && len(file) > 0 {
				fvalue.Set(reflect.ValueOf(file[0]))
			}
		}
	}
	return nil
}

func getFiles(files map[string][]*multipart.FileHeader, names ...string) []*multipart.FileHeader {
	for _, name := range names {
		file, ok := files[name]
		if ok {
			return file
		}
	}
	return nil
}
