package repository

import (
	"ozz-ms/pkg/data/model"
)

func (r Repository) Authorize(username, password string, data interface{}) error {

	if err := r.db.Model(&model.User{}).Where(&model.User{Username: username}).First(data).Error; err != nil {
		return err
	}
	return nil
}
