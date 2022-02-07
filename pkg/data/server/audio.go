package server

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"ozz-ms/pkg/data/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Server) searchAudioRecords(ctx echo.Context) error {

	var err error
	var sp model.AudioRecordingsSearchParams
	if err = ctx.Bind(&sp); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = ctx.Validate(sp); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var data []model.AudioRecording

	if err = s.repo.AudioRecordings(sp, &data); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ret := []model.AudioRecordingDTO{}

	for _, ar := range data {
		ret = append(ret, model.AudioRecordingDTO{
			ID:       ar.ID,
			Name:     ar.Name,
			Path:     ar.Path,
			Duration: ar.Duration,
			Category: ar.Category.Name,
			Date:     ar.Date,
		})
	}

	return ctx.JSON(http.StatusOK, ret)

}

func (s *Server) deleteAudioRecord(ctx echo.Context) error {

	cid := ctx.Param("id")
	if cid == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "No id")
	}
	id, err := strconv.Atoi(cid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ar := model.AudioRecording{}
	if err := s.repo.DeleteAudioRecording(id, &ar); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	// remove files or not?
	//absPath, err := filepath.Abs(ar.Path)
	//if err != nil {
	//	return err
	//}
	//
	//if err = os.Remove(absPath); err != nil {
	//	// do nothing, maybe log warning, as this should not happen
	//}

	return ctx.NoContent(http.StatusOK)
}

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

type ActiveAudioRecordsForCategorySearchParams struct {
	Id   int    `validate:"required|int" param:"id"`
	Name string `query:"name"`
}

func (s *Server) getActiveAudioRecordingsForCategory(ctx echo.Context) error {

	var err error

	sp := ActiveAudioRecordsForCategorySearchParams{}

	if err := ctx.Bind(&sp); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	data := []model.AudioRecording{}
	if err = s.repo.ActiveAudioRecordingsForCategory(sp.Id, sp.Name, &data); err != nil {
		return err
	}

	res := []model.AudioRecordingDTO{}
	for _, ar := range data {
		res = append(res, ar.Map())
	}

	return ctx.JSON(http.StatusOK, res)
}
