package repository

import (
	"errors"

	"ozz-ms/pkg/data/model"

	"gorm.io/gorm"
)

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
