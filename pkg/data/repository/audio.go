package repository

import (
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

func (r Repository) AudioRecording(id int, data interface{}) error {
	return r.db.Model(&model.AudioRecording{}).First(data, id).Error
}
