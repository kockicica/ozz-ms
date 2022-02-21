package repository

import (
	"time"

	"ozz-ms/pkg/data/model"
	"ozz-ms/pkg/util"

	"gorm.io/gorm"
)

func (r Repository) Schedules(sp model.ScheduleSearchParams, data interface{}) error {
	var err error

	tx := r.db.Preload("Recording").Preload("Recording.Category").Preload("Dispositions").Preload("Dispositions.Recording").Preload("Dispositions.Recording.Category")

	//if sp.Category != nil {
	//	tx = tx.Where(&Schedule{Recording: AudioRecording{CategoryID: *sp.Category}})
	//}
	if sp.Recording != nil {
		tx = tx.Where(&model.Schedule{RecordingID: *sp.Recording})
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
	return r.db.
		Preload("Recording").
		Preload("Recording.Category").
		Preload("Dispositions").
		Preload("Dispositions.Recording").
		Preload("Dispositions.Recording.Category").
		First(data, id).Error
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
	sch.TotalPlayCount = data.TotalPlayCount

	columnsToOmit := []string{"TotalPlayCount", "Shift1Played", "Shift2Played", "Shift3Played", "Shift4Played", "Recording", "RecordingID", "Duration"}
	if err := r.db.Select("*").Omit(columnsToOmit...).Updates(&sch).Error; err != nil {
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

func (r Repository) ActiveSchedules(data interface{}) error {
	//act := true
	//TODO: proper filtering of schedules
	return r.Schedules(model.ScheduleSearchParams{}, data)
}
