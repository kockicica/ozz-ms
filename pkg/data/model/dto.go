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
	Client   string
	Comment  string
	Active   bool
	Duration time.Duration
	Date     time.Time
}

type AudioRecordingUpdateDTO struct {
	Name     string `validate:"string"`
	Category string `validate:"string"`
	Client   string `validate:"string"`
	Comment  string `validate:"string"`
	Active   bool   `validate:"bool"`
}

type PagedResults struct {
	Count int64 `json:"count"`
}

type AudioRecordingsPagedResults struct {
	PagedResults
	Data []AudioRecordingDTO `json:"data"`
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
	Shift1 int `validate:"int|min:0"`
	Shift2 int `validate:"int|min:0"`
	Shift3 int `validate:"int|min:0"`
	Shift4 int `validate:"int|min:0"`
	//Active         bool `validate:"bool"`
	TotalPlayCount int
}

type ScheduleDTO struct {
	ID                                                     uint
	Recording                                              AudioRecordingDTO
	Date                                                   time.Time
	Duration                                               time.Duration
	Shift1, Shift2, Shift3, Shift4                         int
	Shift1Played, Shift2Played, Shift3Played, Shift4Played int
	//Active                         bool
	TotalPlayCount int
	Dispositions   []DispositionDTO
}

type AudioRecordingsSearchParams struct {
	Category *int    `query:"category" validate:"int"`
	FromDate *string `validate:"date" query:"fromDate"`
	ToDate   *string `validate:"date" query:"toDate"`
	Active   *bool   `validate:"bool" query:"active"`
	Name     *string `validate:"string" query:"name"`
	Sort     *string `validate:"string" query:"sort"`
	Skip     *int    `validate:"int" query:"skip"`
	Count    *int    `validate:"int" query:"count"`
}

type ScheduleSearchParams struct {
	Recording *int `validate:"int" query:"recording"`
	//Category  *int    `validate:"int" query:"category"`
	//Active   *bool   `validate:"bool" query:"active"`
	FromDate *string `validate:"date" query:"fromDate"`
	ToDate   *string `validate:"date" query:"toDate"`
}

type EqualizerDTO struct {
	ID                                                          uint
	Name                                                        string  `validate:"required"`
	PreAmp                                                      float32 `validate:"float"`
	Amp1, Amp2, Amp3, Amp4, Amp5, Amp6, Amp7, Amp8, Amp9, Amp10 float32 `validate:"float"`
}
