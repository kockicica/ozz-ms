package data

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
	Name       string        `gorm:"unique" validate:"required"`
	Path       string        `gorm:"unique" validate:"required"`
	Duration   time.Duration `validate:"required"`
	Client     *string       `validate:"-"`
	Comment    *string       `validate:"-"`
	CategoryID int
	Category   Category
	Date       time.Time
}

type Disposition struct {
	gorm.Model
	Date             time.Time `validate:"required"`
	Shift            int       `validate:"required"`
	RecordingID      int
	Recording        AudioRecording
	PlayCountNeeded  int `validate:"required|int|min:1|max:1000"`
	PlayCountCurrent int `validate:"required|int|min:1|max:1000"`
}

type DispositionPlayed struct {
	gorm.Model
	DispositionID int
	Disposition   Disposition
	Comment       *string
}

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
	RecordingID int
	Recording   AudioRecording `gorm:"auto_preload:true"`
	Date        time.Time
	Duration    time.Duration
	Shift1      int
	Shift2      int
	Shift3      int
	Shift4      int

	Active         bool
	TotalPlayCount int
}

func (s Schedule) Map() ScheduleDTO {

	return ScheduleDTO{
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
		Active:         s.Active,
		TotalPlayCount: s.TotalPlayCount,
	}

}
