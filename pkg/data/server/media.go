package server

import (
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"ozz-ms/pkg/data/model"

	"github.com/labstack/echo/v4"
)

type AudioRecordingCreateData struct {
	Name     *string               `form:"name" validate:"required"`
	Client   *string               `form:"name"`
	Comment  *string               `form:"comment"`
	Category *string               `form:"category" validate:"required"`
	Duration *string               `form:"duration" validate:"required"`
	File     *multipart.FileHeader `form:"file" validate:"required"`
}

func (s *Server) createAudioRecord(ctx echo.Context) error {

	var err error

	cd := AudioRecordingCreateData{}
	if err = ctx.Bind(&cd); err != nil {
		return err
	}

	if err = ctx.Validate(&cd); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	name := ctx.FormValue("name")
	client := ctx.FormValue("client")
	comment := ctx.FormValue("comment")
	category := ctx.FormValue("category")
	duration := ctx.FormValue("duration")

	file, err := ctx.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get source file
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	// find matching category or default category
	cat, err := s.repo.CategoryByName(category)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	destinationFileName := s.getAudioRecordingPath(file.Filename, cat.Path)

	// create destination folder structure
	destinationFolder := filepath.Dir(destinationFileName)
	if err := os.MkdirAll(destinationFolder, fs.ModeDir); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// create destination file
	dest, err := os.Create(destinationFileName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if _, err = io.Copy(dest, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	_ = dest.Close()

	// create audio recording db record

	// create duration
	dur, err := time.ParseDuration(duration)
	if err != nil {
		// unable to parse duration, this should be treated seriously
		_ = os.Remove(destinationFileName)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ar := model.AudioRecording{
		Name:     name,
		Category: *cat,
		Client:   &client,
		Comment:  &comment,
		Duration: dur,
		Path:     destinationFileName,
		Date:     time.Now(),
	}

	if err := s.repo.NewAudioRecording(&ar); err != nil {
		_ = os.Remove(destinationFileName)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, model.AudioRecordingDTO{
		ID:       ar.ID,
		Name:     ar.Name,
		Category: ar.Category.Name,
		Duration: ar.Duration,
		Path:     ar.Path,
		Date:     ar.Date,
	})
}

//func (s *Server) getCategoryByName(category string) Category {
//	var cat Category
//
//	if err := s.db.Where(&Category{Name: category}).First(&cat).Error; err != nil {
//		if err := s.db.Where(&Category{Default: true}).Find(&cat); err != nil {
//			return Category{Name: "DEFAULT", Path: "default"}
//		}
//	}
//
//	return cat
//}

func (s *Server) getAudioRecordingPath(name, category string) string {
	ext := filepath.Ext(name)
	fileNameWoutExt := name[:len(name)-len(ext)]
	cd := time.Now().Format("20060102150405")

	return filepath.Join(s.Config.RootPath, category, fmt.Sprintf("%s-%s%s", fileNameWoutExt, cd, ext))
}
