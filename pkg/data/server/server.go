package server

import (
	"context"
	"fmt"
	"mime/multipart"
	"reflect"
	"time"

	"ozz-ms/pkg/data/repository"

	"github.com/labstack/echo/v4"
)

var (
	typeMultipartFileHeader      = reflect.TypeOf((*multipart.FileHeader)(nil)).Elem()
	typeMultipartSliceFileHeader = reflect.TypeOf(([]*multipart.FileHeader)(nil)).Elem()
)

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

	repoCfg := repository.RepositoryConfig{
		Dsn:     s.Config.Dsn,
		Verbose: s.Config.Verbose,
	}
	r, err := repository.NewRepository(repoCfg)
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
	audioGroup.POST("", ds.createAudioRecord)
	audioGroup.GET("", ds.searchAudioRecords)
	audioGroup.PUT("/:id", ds.updateAudioRecord)
	audioGroup.DELETE("/:id", ds.deleteAudioRecord)
	audioGroup.GET("/media/:id", ds.serveAudioFile)
	//audioGroup.GET("/active/:id", ds.getActiveAudioRecordingsForCategory)

	scheduleGroup := apiGroup.Group("/schedules")
	scheduleGroup.GET("", ds.searchSchedules)
	scheduleGroup.GET("/:id", ds.getSchedule)
	scheduleGroup.PUT("/:id", ds.updateSchedule)
	scheduleGroup.DELETE("/:id", ds.deleteSchedule)
	scheduleGroup.POST("", ds.createSchedule)
	scheduleGroup.POST("/multiple", ds.createMultipleSchedules)

	dispositionGroup := apiGroup.Group("/dispositions")
	dispositionGroup.GET("", ds.searchDispositions)
	dispositionGroup.POST("/create", ds.createDispositions)
	dispositionGroup.POST("/mark", ds.markDispositionExecution)
	//dispositionGroup.POST("/:id/increase", ds.increaseDispositionPlayedCount)
	//dispositionGroup.POST("/:id/decrease", ds.decreaseDispositionPlayedCount)

	equalizerGroup := apiGroup.Group("/equalizers")
	equalizerGroup.GET("", ds.getEqualizers)
	equalizerGroup.GET("/:id", ds.getEqualizer)
	equalizerGroup.POST("", ds.createEqualizer)
	equalizerGroup.PUT("/:id", ds.updateEqualizer)
	equalizerGroup.DELETE("/:id", ds.deleteEqualizer)

	apiGroup.GET("/shifts", ds.getShifts)
	apiGroup.POST("/authorize", ds.authorize)
	apiGroup.GET("/categories", ds.getCategories)

	return ds, nil
}
