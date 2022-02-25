package repository

import (
	"time"

	"ozz-ms/pkg/data/model"
)

//func (r Repository) CreateDispositions(sch *model.Schedule) error {
//
//	dispositionsToCreate := []model.Disposition{}
//
//	var disp model.Disposition
//
//	disp = model.Disposition{
//		PlayCountCurrent: 0,
//		PlayCountNeeded:  sch.Shift1,
//		Date:             sch.Date,
//		Shift:            1,
//		RecordingID:      sch.RecordingID,
//		Recording:        sch.Recording,
//		Schedule:         *sch,
//	}
//	dispositionsToCreate = append(dispositionsToCreate, disp)
//
//	disp = model.Disposition{
//		PlayCountCurrent: 0,
//		PlayCountNeeded:  sch.Shift2,
//		Date:             sch.Date,
//		Shift:            2,
//		RecordingID:      sch.RecordingID,
//		Recording:        sch.Recording,
//		Schedule:         *sch,
//	}
//	dispositionsToCreate = append(dispositionsToCreate, disp)
//
//	disp = model.Disposition{
//		PlayCountCurrent: 0,
//		PlayCountNeeded:  sch.Shift3,
//		Date:             sch.Date,
//		Shift:            3,
//		RecordingID:      sch.RecordingID,
//		Recording:        sch.Recording,
//		Schedule:         *sch,
//	}
//	dispositionsToCreate = append(dispositionsToCreate, disp)
//
//	disp = model.Disposition{
//		PlayCountCurrent: 0,
//		PlayCountNeeded:  sch.Shift4,
//		Date:             sch.Date,
//		Shift:            4,
//		RecordingID:      sch.RecordingID,
//		Recording:        sch.Recording,
//		Schedule:         *sch,
//	}
//	dispositionsToCreate = append(dispositionsToCreate, disp)
//
//	// find current child dispositions
//	currentDispositions := []model.Disposition{}
//	for _, cd := range sch.Dispositions {
//		currentDispositions = append(currentDispositions, cd)
//	}
//
//	if err := r.db.Model(sch).Association("Dispositions").Replace(dispositionsToCreate); err != nil {
//		return err
//	}
//
//	if len(currentDispositions) > 0 {
//		if err := r.db.Unscoped().Delete(&currentDispositions).Error; err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//

func (r Repository) DispositionForShiftAndData(shift int, date time.Time) ([]model.DispositionDTO, error) {

	data := []model.DispositionDTO{}

	// find schedule data
	schedules := []model.Schedule{}

	tx := r.db.Preload("Recording").
		Preload("Recording.Category").
		Where("Date = ?", date)
	if err := tx.Find(&schedules).Error; err != nil {
		return nil, err
	}

	for _, schedule := range schedules {
		needed := 0
		played := 0
		switch shift {
		case 1:
			needed = schedule.Shift1 + (schedule.Shift4 - schedule.Shift4Played)
			played = schedule.Shift1Played
		case 2:
			needed = schedule.Shift2 + (schedule.Shift1 - schedule.Shift1Played) + +(schedule.Shift4 - schedule.Shift4Played)
			played = schedule.Shift2Played
		case 3:
			needed = schedule.Shift3 + (schedule.Shift1 - schedule.Shift1Played) + (schedule.Shift2 - schedule.Shift2Played) + (schedule.Shift4 - schedule.Shift4Played)
			played = schedule.Shift3Played
		case 4:
			needed = schedule.Shift4
			played = schedule.Shift4Played
		}
		data = append(data, model.DispositionDTO{
			AudioRecordingDTO: schedule.Recording.Map(),
			Date:              schedule.Date,
			Shift:             shift,
			PlayCountNeeded:   needed,
			PlayCountCurrent:  played,
			ScheduleID:        schedule.ID,
		})
	}

	type PrevResult struct {
		Extra int
	}
	var extra PrevResult
	// second pass, adding missed
	for i, disposition := range data {
		tx := r.db.Model(&model.Schedule{}).
			Where("Date < ? and Recording_ID = ?", date, disposition.AudioRecordingDTO.ID).
			Select("sum(shift1+shift2+shift3+shift4) - sum(shift1_played+shift2_played+shift3_played+shift4_played) as extra").
			Scan(&extra)
		if tx.Error != nil {
			return nil, tx.Error
		}
		data[i].PlayCountNeeded += extra.Extra
		data[i].PlayCountRemaining = data[i].PlayCountNeeded - data[i].PlayCountCurrent
	}

	return data, nil
}

//
//func (r Repository) ChangePlayCountForDisposition(id int, time time.Time, delta int) (*model.Disposition, error) {
//
//	tmp := model.Disposition{}
//
//	if erro := r.db.Transaction(func(tx *gorm.DB) error {
//		if err := tx.Preload("Recording").First(&tmp, id).Error; err != nil {
//			return err
//		}
//		if delta > 0 {
//			if tmp.PlayCountCurrent < tmp.PlayCountNeeded {
//				tmp.PlayCountCurrent += delta
//			}
//		} else {
//			if tmp.PlayCountCurrent > 0 {
//				tmp.PlayCountCurrent += delta
//			}
//		}
//		if err := tx.Select("PlayCountCurrent").Updates(&tmp).Error; err != nil {
//			return err
//		}
//
//		return nil
//	}); erro != nil {
//		return nil, erro
//	}
//
//	return &tmp, nil
//
//}

func (r Repository) CreateDispositions(cdp model.CreateDispositionParams, data *[]model.Schedule) error {

	// find schedules
	fromDate, err := time.Parse("2006-01-02", cdp.From)
	if err != nil {
		return err
	}

	toDate := fromDate.AddDate(0, 0, cdp.Days)

	local := []model.Schedule{}

	if err := r.db.Model(&model.Schedule{}).Where("Date >= ? and Date < ?", fromDate, toDate).Find(&local).Error; err != nil {
		return err
	}

	for _, schedule := range local {
		schedule.HasDisposition = true
		if err := r.db.Model(&schedule).Select("HasDisposition").Updates(&schedule).Error; err != nil {
			return err
		}
	}

	if err := r.db.Model(&model.Schedule{}).Where("Date >= ? and Date < ?", fromDate, toDate).Find(data).Error; err != nil {
		return err
	}

	return nil
}

func (r Repository) MarkDispositionExecute(data model.DispositionExecuteParams) error {

	var schedule model.Schedule

	if err := r.db.Model(&model.Schedule{}).First(&schedule, data.Schedule).Error; err != nil {
		return err
	}

	switch data.Shift {
	case 1:
		schedule.Shift1Played += 1
	case 2:
		schedule.Shift2Played += 1
	case 3:
		schedule.Shift3Played += 1
	case 4:
		schedule.Shift4Played += 1
	}

	if err := r.db.Model(&schedule).Select("Shift1Played", "Shift2Played", "Shift3Played", "Shift4Played").Updates(&schedule).Error; err != nil {
		return err
	}

	return nil
}
