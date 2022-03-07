package repository

import (
	"fmt"
	"strings"
	"time"

	"ozz-ms/pkg/data/model"
)

func (r Repository) NewAudioRecording(rec *model.AudioRecording) error {
	if err := r.db.Create(rec).Error; err != nil {
		return err
	}
	return nil
}

func (r Repository) AudioRecordings(sp model.AudioRecordingsSearchParams, data interface{}, count *int64) error {
	tx := r.db.Preload("Category").Model(&model.AudioRecording{})

	if sp.Category != nil {
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

	if sp.Active != nil {
		tx = tx.Where("Active = ?", *sp.Active)
	}

	if sp.Name != nil {
		tx = tx.Where("Name like ?", fmt.Sprintf("%%%s%%", *sp.Name))
	}

	if sp.Sort != nil {
		var sortClause string
		desc := strings.HasPrefix(*sp.Sort, "-")
		if desc {
			sortClause = fmt.Sprintf("%s desc", (*sp.Sort)[1:])
		} else {
			sortClause = *sp.Sort
		}
		tx = tx.Order(sortClause)
	}

	tx = tx.Count(count)

	if sp.Skip != nil {
		tx = tx.Offset(*sp.Skip)
	} else {
		tx = tx.Offset(0)
	}
	if sp.Count != nil {
		tx = tx.Limit(*sp.Count)
	} else {
		tx = tx.Limit(20)
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

	if err := r.db.Unscoped().Delete(&model.AudioRecording{}, id).Error; err != nil {
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

func (r Repository) UpdateAudioRecording(id int, updateData *model.AudioRecordingUpdateDTO, data interface{}) error {

	fnd := model.AudioRecording{}

	if err := r.db.Model(&model.AudioRecording{}).Preload("Category").First(&fnd, id).Error; err != nil {
		return err
	}

	cat := model.Category{}
	if err := r.db.Model(&model.Category{}).Where("Name = ?", updateData.Category).First(&cat).Error; err != nil {
		return err
	}

	updateDict := map[string]interface{}{
		"Name":     updateData.Name,
		"Client":   updateData.Client,
		"Comment":  updateData.Comment,
		"Active":   updateData.Active,
		"Category": cat,
	}

	if err := r.db.Model(&fnd).Updates(updateDict).Error; err != nil {
		return err
	}

	if err := r.db.Preload("Category").First(data, id).Error; err != nil {
		return err
	}

	return nil
}

func (r Repository) GetAudioLog(sp model.AudioRecordingLogSearchParams, data interface{}) error {

	tx := r.db.
		Joins("Schedule").
		Preload("Schedule.Recording").
		Preload("Schedule.Recording.Category").
		Model(&model.EmitLog{}).
		Where("Schedule.Recording_ID = ?", sp.Recording)

	if sp.From != nil {
		dateFrom, err := time.Parse("2006-01-02", *sp.From)
		if err != nil {
			return err
		}
		tx = tx.Where("Time >= ?", dateFrom)
	}

	if sp.To != nil {
		dateTo, err := time.Parse("2006-01-02", *sp.To)
		if err != nil {
			return err
		}
		tx = tx.Where("Time <= ?", dateTo)
	}

	tx = tx.Order("Time")

	if err := tx.Find(data).Error; err != nil {
		return err
	}

	return nil
}
