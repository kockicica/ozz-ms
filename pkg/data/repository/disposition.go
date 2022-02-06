package repository

import (
	"ozz-ms/pkg/data/model"
)

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
