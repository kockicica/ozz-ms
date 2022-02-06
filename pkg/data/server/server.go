package server

import (
	"context"
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"
	"time"

	"ozz-ms/pkg/data/repository"

	"github.com/gookit/validate"
	"github.com/labstack/echo/v4"
)

var (
	typeMultipartFileHeader      = reflect.TypeOf((*multipart.FileHeader)(nil)).Elem()
	typeMultipartSliceFileHeader = reflect.TypeOf(([]*multipart.FileHeader)(nil)).Elem()
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

type ServerConfig struct {
	Port     int
	Dsn      string
	Verbose  bool
	RootPath string
}

type Server struct {
	Config ServerConfig
	es     *echo.Echo
	//db     *gorm.DB
	repo *repository.Repository
}

func (s *Server) Start() error {
	var err error

	r, err := repository.NewSQLiteRepository(s.Config.Dsn)
	if err != nil {
		return err
	}

	s.repo = r

	err = s.es.Start(fmt.Sprintf(":%d", s.Config.Port))
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.es.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewDataServer(config ServerConfig) (*Server, error) {

	ds := new(Server)
	ds.Config = config
	ds.es = echo.New()
	ds.es.HideBanner = true
	ds.es.HidePort = true
	ds.es.Validator = NewServerValidator()
	ds.es.Binder = &ServerBinding{}

	apiGroup := ds.es.Group("/api")

	audioGroup := apiGroup.Group("/audio")
	audioGroup.POST("", ds.uploadAudioTrack)
	audioGroup.GET("", ds.searchAudioRecords)
	audioGroup.DELETE("/:id", ds.deleteAudioRecord)

	scheduleGroup := apiGroup.Group("/schedules")
	scheduleGroup.GET("", ds.searchSchedules)
	scheduleGroup.GET("/:id", ds.getSchedule)
	scheduleGroup.PUT("/:id", ds.updateSchedule)
	scheduleGroup.DELETE("/:id", ds.deleteSchedule)
	scheduleGroup.POST("", ds.createSchedule)

	apiGroup.GET("/shifts", ds.getShifts)
	apiGroup.POST("/authorize", ds.authorize)
	apiGroup.GET("/categories", ds.getCategories)

	return ds, nil
}

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
