package model

import (
	"time"
)

type CategoryDTO struct {
	ID    int
	Name  string
	Order int
	Path  string
}

type AudioRecordingDTO struct {
	ID       uint
	Name     string
	Path     string
	Category string
	Duration time.Duration
	Date     time.Time
}

type DispositionDTO struct {
	AudioRecordingDTO
	Date             time.Time
	Shift            int
	PlayCountNeeded  int
	PlayCountCurrent int
	RecordingID      int
}

type NewScheduleDTO struct {
	ID        uint
	Recording int    `validate:"required|int"`
	Date      string `validate:"required|date"`
	//Duration       string `validate:"required"`
	Shift1         int  `validate:"int|min:0"`
	Shift2         int  `validate:"int|min:0"`
	Shift3         int  `validate:"int|min:0"`
	Shift4         int  `validate:"int|min:0"`
	Active         bool `validate:"bool"`
	TotalPlayCount int
}

type ScheduleDTO struct {
	ID                             uint
	Recording                      AudioRecordingDTO
	Date                           time.Time
	Duration                       time.Duration
	Shift1, Shift2, Shift3, Shift4 int
	Active                         bool
	TotalPlayCount                 int
	Dispositions                   []DispositionDTO
}

type AudioRecordingsSearchParams struct {
	Category *int    `query:"category" validate:"int"`
	FromDate *string `validate:"date" query:"fromDate"`
	ToDate   *string `validate:"date" query:"toDate"`
}

type ScheduleSearchParams struct {
	Recording *int `validate:"int" query:"recording"`
	//Category  *int    `validate:"int" query:"category"`
	Active   *bool   `validate:"bool" query:"active"`
	FromDate *string `validate:"date" query:"fromDate"`
	ToDate   *string `validate:"date" query:"toDate"`
}
