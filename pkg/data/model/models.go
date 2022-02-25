package model

import (
	"fmt"
	"time"

	"github.com/gosuri/uitable"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Shift struct {
	gorm.Model
	Name   string `gorm:"unique" validate:"required" message:"Name is required"`
	Order  int
	Active bool `validate:"-"`
}

type Shifts []Shift

func (s Shifts) Print() {
	table := uitable.New()
	table.AddRow("ID")
	table.AddRow("NAME")
	table.AddRow("ACTIVE")
	for _, v := range s {
		table.AddRow(v.ID)
		table.AddRow(v.Name)
		table.AddRow(v.Active)
	}
	fmt.Println(table)
}

type Category struct {
	gorm.Model
	Name    string `validate:"required"`
	Order   int    `validate:"-"`
	Path    string `validate:"-"`
	Default bool
}

type AudioRecording struct {
	gorm.Model
	Name       string        `validate:"required" gorm:"index:one_name,unique,where:deleted_at is null"`
	Path       string        `validate:"required"`
	Duration   time.Duration `validate:"required"`
	Client     *string       `validate:"-"`
	Comment    *string       `validate:"-"`
	Active     bool
	CategoryID int
	Category   Category
	Date       time.Time
}

func (r AudioRecording) Map() AudioRecordingDTO {
	return AudioRecordingDTO{
		ID:       r.ID,
		Name:     r.Name,
		Path:     r.Path,
		Category: r.Category.Name,
		Client:   *r.Client,
		Comment:  *r.Comment,
		Active:   r.Active,
		Duration: r.Duration,
		Date:     r.Date,
	}
}

//type Disposition struct {
//	gorm.Model
//	Date             time.Time `validate:"required"`
//	Shift            int       `validate:"required"`
//	RecordingID      int
//	Recording        AudioRecording
//	PlayCountNeeded  int `validate:"required|int|min:1|max:1000"`
//	PlayCountCurrent int `validate:"required|int|min:1|max:1000"`
//	ScheduleID       int
//	Schedule         Schedule
//}
//
//func (d *Disposition) Map() DispositionDTO {
//	dto := DispositionDTO{
//		AudioRecordingDTO: d.Recording.Map(),
//		Date:              d.Date,
//		Shift:             d.Shift,
//		PlayCountNeeded:   d.PlayCountNeeded,
//		PlayCountCurrent:  d.PlayCountCurrent,
//		RecordingID:       d.RecordingID,
//	}
//	dto.ID = d.ID
//	dto.Date = d.Date
//	return dto
//}
//
//type DispositionPlayed struct {
//	gorm.Model
//	DispositionID int
//	Disposition   Disposition
//	Comment       *string
//}

type UserLevel int

const (
	Regular UserLevel = iota
	Admin
)

type User struct {
	gorm.Model
	Username string
	Password string
	Level    UserLevel
}

func (u *User) SetPassword(password string) error {

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)

	return nil
}

type Schedule struct {
	gorm.Model
	RecordingID  int
	Recording    AudioRecording `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Date         time.Time
	Duration     time.Duration
	Shift1       int
	Shift2       int
	Shift3       int
	Shift4       int
	Shift1Played int
	Shift2Played int
	Shift3Played int
	Shift4Played int

	//Active         bool
	TotalPlayCount int

	HasDisposition bool
	//Dispositions []Disposition
}

func (s Schedule) Map() ScheduleDTO {

	dto := ScheduleDTO{
		ID: s.ID,
		Recording: AudioRecordingDTO{
			ID:       s.Recording.ID,
			Name:     s.Recording.Name,
			Path:     s.Recording.Path,
			Category: s.Recording.Category.Name,
			Duration: s.Recording.Duration,
			Date:     s.Recording.Date,
		},
		Date:           s.Date,
		Duration:       s.Duration,
		Shift1:         s.Shift1,
		Shift2:         s.Shift2,
		Shift3:         s.Shift3,
		Shift4:         s.Shift4,
		Shift1Played:   s.Shift1Played,
		Shift2Played:   s.Shift2Played,
		Shift3Played:   s.Shift3Played,
		Shift4Played:   s.Shift4Played,
		TotalPlayCount: s.TotalPlayCount,
		HasDisposition: s.HasDisposition,
		//Dispositions:   []DispositionDTO{},
	}

	//for _, dsp := range s.Dispositions {
	//	dto.Dispositions = append(dto.Dispositions, dsp.Map())
	//}

	return dto

}

type Equalizer struct {
	gorm.Model
	Name                                                        string `gorm:"unique"`
	PreAmp                                                      float32
	Amp1, Amp2, Amp3, Amp4, Amp5, Amp6, Amp7, Amp8, Amp9, Amp10 float32
}

func (e Equalizer) Map() EqualizerDTO {
	return EqualizerDTO{
		ID:     e.ID,
		Name:   e.Name,
		PreAmp: e.PreAmp,
		Amp1:   e.Amp1,
		Amp2:   e.Amp2,
		Amp3:   e.Amp3,
		Amp4:   e.Amp4,
		Amp5:   e.Amp5,
		Amp6:   e.Amp6,
		Amp7:   e.Amp7,
		Amp8:   e.Amp8,
		Amp9:   e.Amp9,
		Amp10:  e.Amp10,
	}
}
