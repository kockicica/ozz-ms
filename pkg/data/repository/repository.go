package repository

import (
	"errors"
	"time"

	"ozz-ms/pkg/data/model"
	"ozz-ms/pkg/util"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	db *gorm.DB
}

func (r Repository) Authorize(username, password string, data interface{}) error {

	if err := r.db.Model(&model.User{}).Where(&model.User{Username: username}).First(data).Error; err != nil {
		return err
	}
	return nil
}
func (r Repository) Shifts(data interface{}) error {

	if err := r.db.Model(&model.Shift{}).Order("`order`").Find(data).Error; err != nil {
		return err
	}

	return nil
}

func (r Repository) Categories(data interface{}) error {

	if err := r.db.Model(&model.Category{}).Order("`order`").Find(data).Error; err != nil {
		return err
	}

	return nil
}

func (r Repository) CategoryByName(name string) (*model.Category, error) {
	var cat = new(model.Category)
	if err := r.db.Where(&model.Category{Name: name}).Find(cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// not found, make default
			return &model.Category{Name: "DEFAULT", Path: "default"}, nil
		}
		return nil, err
	}
	return cat, nil
}

func (r Repository) NewAudioRecording(rec *model.AudioRecording) error {
	if err := r.db.Create(rec).Error; err != nil {
		return err
	}
	return nil
}

type ScheduleSearchParams struct {
	Recording *int `validate:"int" query:"recording"`
	//Category  *int    `validate:"int" query:"category"`
	Active   *bool   `validate:"bool" query:"active"`
	FromDate *string `validate:"date" query:"fromDate"`
	ToDate   *string `validate:"date" query:"toDate"`
}

func (r Repository) Schedules(sp ScheduleSearchParams, data interface{}) error {
	var err error

	tx := r.db.Preload("Recording").Preload("Recording.Category")

	//if sp.Category != nil {
	//	tx = tx.Where(&Schedule{Recording: AudioRecording{CategoryID: *sp.Category}})
	//}
	if sp.Recording != nil {
		tx = tx.Where(&model.Schedule{RecordingID: *sp.Recording})
	}
	if sp.Active != nil {
		tx = tx.Where("Active", *sp.Active)
	}
	if sp.FromDate != nil {
		fdt, err := time.Parse("2006-01-02", *sp.FromDate)
		if err != nil {
			return err
		}
		tx = tx.Where("Date >= ?", fdt)
	}
	if sp.ToDate != nil {
		fdt, err := time.Parse("2006-01-02", *sp.ToDate)
		if err != nil {
			return err
		}
		tx = tx.Where("Date <= ?", fdt)
	}

	if err = tx.Find(data).Error; err != nil {
		return err
	}
	return nil

}

func (r Repository) Schedule(id int, data interface{}) error {
	return r.db.Preload("Recording").Preload("Recording.Category").First(data, id).Error
}

func (r Repository) DeleteSchedule(id []int) error {
	tx := r.db.Delete(&model.Schedule{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil

}

func (r Repository) SetSchedule(id int, data model.NewScheduleDTO) error {

	sch := &model.Schedule{}

	tx := r.db.Preload("Recording").First(&sch, id)
	if tx.Error != nil {
		return tx.Error
	}

	scheduleDate, err := util.ParseDateString(data.Date)
	if err != nil {
		return err
	}

	sch.RecordingID = data.Recording
	sch.Date = scheduleDate
	sch.Duration = sch.Recording.Duration
	sch.Shift1 = data.Shift1
	sch.Shift2 = data.Shift2
	sch.Shift3 = data.Shift3
	sch.Shift4 = data.Shift4
	sch.Active = data.Active
	sch.TotalPlayCount = data.TotalPlayCount

	if err := r.db.Select("*").Omit("TotalPlayCount").Updates(&sch).Error; err != nil {
		return err
	}
	return nil
}

type AudioRecordingsSearchParams struct {
	Category *int    `query:"category" validate:"int"`
	FromDate *string `validate:"date" query:"fromDate"`
	ToDate   *string `validate:"date" query:"toDate"`
}

func (r Repository) AudioRecordings(sp AudioRecordingsSearchParams, data interface{}) error {
	tx := r.db.Preload("Category").Model(&model.AudioRecording{})

	if sp.Category != nil {
		tx = tx.Where("Category.id", *sp.Category)
	}

	if sp.FromDate != nil {
		fdt, err := time.Parse("2006-01-02", *sp.FromDate)
		if err != nil {
			return err
		}
		tx = tx.Where("Date >= ?", fdt)
	}

	if sp.ToDate != nil {
		fdt, err := time.Parse("2006-01-02", *sp.ToDate)
		if err != nil {
			return err
		}
		tx = tx.Where("Date <= ?", fdt)
	}
	if err := tx.Find(data).Error; err != nil {
		return err
	}
	return nil
}

func (r Repository) DeleteAudioRecording(id int, data interface{}) error {

	if err := r.db.Model(&model.AudioRecording{}).First(data, id).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&model.AudioRecording{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r Repository) NewSchedule(dto model.NewScheduleDTO) (*model.Schedule, error) {

	//d, err := time.ParseDuration(dto.Duration)
	//if err != nil {
	//	return nil, err
	//}

	dd, err := time.Parse("2006-01-02", dto.Date)
	if err != nil {
		return nil, err
	}

	// find
	rec := model.AudioRecording{}
	if err := r.db.First(&rec, dto.Recording).Error; err != nil {
		return nil, err
	}

	sch := model.Schedule{
		Duration:       rec.Duration,
		Shift1:         dto.Shift1,
		Shift2:         dto.Shift2,
		Shift3:         dto.Shift3,
		Shift4:         dto.Shift4,
		Active:         true,
		Date:           dd,
		TotalPlayCount: 0,
		RecordingID:    dto.Recording,
	}

	if err := r.db.Create(&sch).Error; err != nil {
		return nil, err
	}

	if err := r.db.Preload("Recording").Preload("Recording.Category").Find(&sch).Error; err != nil {
		return nil, err
	}

	return &sch, nil

}

func (r Repository) CreateDispositions(items []model.Schedule) ([]model.Disposition, error) {

	for _, sch := range items {
		var disp model.Disposition
		disp = model.Disposition{
			PlayCountCurrent: 0,
			PlayCountNeeded:  sch.Shift1,
		}
		r.db.Create(&disp)
	}

	return nil, nil
}

func NewSQLiteRepository(dsn string) (*Repository, error) {

	var err error
	var db *gorm.DB

	repo := new(Repository)

	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		//SkipDefaultTransaction: true,
	})
	//db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}

	models := []interface{}{
		&model.Shift{},
		&model.Category{},
		&model.AudioRecording{},
		&model.Disposition{},
		&model.DispositionPlayed{},
		&model.User{},
		&model.Schedule{},
	}

	if err = db.AutoMigrate(models...); err != nil {
		return nil, err
	}

	if err = initCategories(db); err != nil {
		return nil, err
	}
	if err = initShifts(db); err != nil {
		return nil, err
	}

	if err = initUsers(db); err != nil {
		return nil, err
	}

	repo.db = db

	return repo, nil

}

func initCategories(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Category{}).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return nil
	}
	var predefinedCategories = []model.Category{
		{
			Name:  "REKLAME",
			Path:  "reklame",
			Order: 1,
		},
		{
			Name:  "ŠPICE",
			Path:  "spice",
			Order: 2,
		},
		{
			Name:  "UPADICE",
			Path:  "upadice",
			Order: 3,
		},
		{
			Name:  "MASKE",
			Path:  "maske",
			Order: 3,
		},
		{
			Name:  "DŽINGLOVI",
			Path:  "dzinglovi",
			Order: 4,
		},
		{
			Name:  "SLOBODNA",
			Path:  "slobodna",
			Order: 5,
		},
	}

	for _, c := range predefinedCategories {
		if err := db.Create(&c).Error; err != nil {
			return err
		}
	}
	return nil
}

func initShifts(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Shift{}).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return nil
	}

	predefinedShifts := []model.Shift{
		{Name: "Smena I", Active: true, Order: 2},
		{Name: "Smena II", Active: true, Order: 3},
		{Name: "Smena III", Active: true, Order: 4},
		{Name: "Smena IV", Active: true, Order: 1},
	}

	for _, s := range predefinedShifts {
		if err := db.Create(&s).Error; err != nil {
			return err
		}
	}

	return nil
}

func initUsers(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return nil
	}

	predefinedUsers := []model.User{
		{Username: "maki", Password: "maki", Level: model.Admin},
		{Username: "taki", Password: "taki", Level: model.Regular},
		{Username: "laki", Password: "laki", Level: model.Regular},
		{Username: "caki", Password: "caki", Level: model.Regular},
	}

	for _, u := range predefinedUsers {
		if err := u.SetPassword(u.Username); err != nil {
			return err
		}
		if err := db.Create(&u).Error; err != nil {
			return err
		}
	}

	return nil
}
