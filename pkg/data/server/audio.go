package server

import (
	"errors"
	"net/http"
	"strconv"

	"ozz-ms/pkg/data/model"
	"ozz-ms/pkg/data/repository"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Server) searchAudioRecords(ctx echo.Context) error {

	var err error
	var sp repository.AudioRecordingsSearchParams
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
