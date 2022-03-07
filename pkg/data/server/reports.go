package server

import (
	"net/http"

	"ozz-ms/pkg/data/model"

	"github.com/labstack/echo/v4"
)

func (s *Server) audioRecordingLog(ctx echo.Context) error {

	var sp model.AudioRecordingLogSearchParams

	if err := ctx.Bind(&sp); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	logs := []model.EmitLog{}
	if err := s.repo.GetAudioLog(sp, &logs); err != nil {
		return err
	}

	report := []model.AudioRecordingLog{}

	for _, log := range logs {
		report = append(report, model.AudioRecordingLog{
			Name:         log.Schedule.Recording.Name,
			Category:     log.Schedule.Recording.Category.Name,
			Duration:     log.Schedule.Recording.Duration,
			Time:         log.Time,
			Shift:        log.Shift,
			ScheduleDate: log.Schedule.Date,
		})
	}

	return ctx.JSON(http.StatusOK, report)
}
