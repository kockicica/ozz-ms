package repository

import (
	"time"

	"ozz-ms/pkg/data/model"

	"gorm.io/gorm"
)

func (r Repository) CreateDispositions(sch *model.Schedule) error {

	dispositionsToCreate := []model.Disposition{}

	var disp model.Disposition

	disp = model.Disposition{
		PlayCountCurrent: 0,
		PlayCountNeeded:  sch.Shift1,
		Date:             sch.Date,
		Shift:            1,
		RecordingID:      sch.RecordingID,
		Recording:        sch.Recording,
		Schedule:         *sch,
	}
	dispositionsToCreate = append(dispositionsToCreate, disp)

	disp = model.Disposition{
		PlayCountCurrent: 0,
		PlayCountNeeded:  sch.Shift2,
		Date:             sch.Date,
		Shift:            2,
		RecordingID:      sch.RecordingID,
		Recording:        sch.Recording,
		Schedule:         *sch,
	}
	dispositionsToCreate = append(dispositionsToCreate, disp)

	disp = model.Disposition{
		PlayCountCurrent: 0,
		PlayCountNeeded:  sch.Shift3,
		Date:             sch.Date,
		Shift:            3,
		RecordingID:      sch.RecordingID,
		Recording:        sch.Recording,
		Schedule:         *sch,
	}
	dispositionsToCreate = append(dispositionsToCreate, disp)

	disp = model.Disposition{
		PlayCountCurrent: 0,
		PlayCountNeeded:  sch.Shift4,
		Date:             sch.Date,
		Shift:            4,
		RecordingID:      sch.RecordingID,
		Recording:        sch.Recording,
		Schedule:         *sch,
	}
	dispositionsToCreate = append(dispositionsToCreate, disp)

	// find current child dispositions
	currentDispositions := []model.Disposition{}
	for _, cd := range sch.Dispositions {
		currentDispositions = append(currentDispositions, cd)
	}

	if err := r.db.Model(sch).Association("Dispositions").Replace(dispositionsToCreate); err != nil {
		return err
	}

	if len(currentDispositions) > 0 {
		if err := r.db.Unscoped().Delete(&currentDispositions).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r Repository) DispositionForShiftAndData(shift int, date time.Time) ([]model.Disposition, error) {

	data := []model.Disposition{}

	if err := r.db.Preload("Recording").
		Preload("Recording.Category").
		Find(&data, &model.Disposition{Shift: shift, Date: date}).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r Repository) ChangePlayCountForDisposition(id int, time time.Time, delta int) (*model.Disposition, error) {

	tmp := model.Disposition{}

	if erro := r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.db.Preload("Recording").First(&tmp, id).Error; err != nil {
			return err
		}
		if delta > 0 {
			if tmp.PlayCountCurrent < tmp.PlayCountNeeded {
				tmp.PlayCountCurrent += delta
			}
		} else {
			if tmp.PlayCountCurrent > 0 {
				tmp.PlayCountCurrent += delta
			}
		}
		if err := r.db.Updates(&tmp).Error; err != nil {
			return err
		}

		return nil
	}); erro != nil {
		return nil, erro
	}

	return &tmp, nil

}
