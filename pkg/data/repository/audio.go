package repository

import (
	"fmt"
	"time"

	"ozz-ms/pkg/data/model"
)

func (r Repository) NewAudioRecording(rec *model.AudioRecording) error {
	if err := r.db.Create(rec).Error; err != nil {
		return err
	}
	return nil
}

func (r Repository) AudioRecordings(sp model.AudioRecordingsSearchParams, data interface{}) error {
	tx := r.db.Preload("Category").Model(&model.AudioRecording{})

	if sp.Category != nil {
		//tx = tx.Where("CategoryID", *sp.Category)
		tx = tx.Where(&model.AudioRecording{CategoryID: *sp.Category})
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

func (r Repository) AudioRecording(id int, data interface{}) error {
	return r.db.Model(&model.AudioRecording{}).First(data, id).Error
}

func (r Repository) ActiveAudioRecordingsForCategory(catId int, name string, data interface{}) error {

	tx := r.db.Model(&model.AudioRecording{}).
		Select("distinct Audio_Recordings.*, Categories.name as Category__name, Categories.id as Category__id, Categories.`Order` as Category__order").
		Joins("join Categories on Categories.Id = Audio_Recordings.Category_ID").
		Joins("left outer join Schedules on Schedules.Recording_id = Audio_Recordings.Id and Schedules.active = ?", true).
		Where("Categories.Id = ?", catId)

	if name != "" {
		tx = tx.Where("Audio_Recordings.Name like ?", fmt.Sprintf("%%%s%%", name))
	}

	if err := tx.Find(data).Error; err != nil {
		return err
	}

	return nil
}
