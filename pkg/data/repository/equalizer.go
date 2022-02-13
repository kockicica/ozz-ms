package repository

import (
	"ozz-ms/pkg/data/model"
)

func (r Repository) Equalizers(data interface{}) error {

	if err := r.db.Model(&model.Equalizer{}).Find(data).Error; err != nil {
		return err
	}

	return nil
}

func (r Repository) Equalizer(id int, data interface{}) error {

	if err := r.db.Model(&model.Equalizer{}).First(data, id).Error; err != nil {
		return err
	}

	return nil
}

func (r Repository) NewEqualizer(data interface{}) error {

	if err := r.db.Model(&model.Equalizer{}).Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (r Repository) SetEqualizer(id int, data model.EqualizerDTO) error {

	fnd := model.Equalizer{}
	if err := r.db.First(&fnd, id).Error; err != nil {
		return err
	}

	fnd.Name = data.Name
	fnd.PreAmp = data.PreAmp
	fnd.Amp1 = data.Amp1
	fnd.Amp2 = data.Amp2
	fnd.Amp3 = data.Amp3
	fnd.Amp3 = data.Amp3
	fnd.Amp4 = data.Amp4
	fnd.Amp5 = data.Amp5
	fnd.Amp6 = data.Amp6
	fnd.Amp7 = data.Amp7
	fnd.Amp8 = data.Amp8
	fnd.Amp9 = data.Amp9
	fnd.Amp10 = data.Amp10

	if err := r.db.Select("*").Updates(&fnd).Error; err != nil {
		return err
	}

	return nil
}

func (r Repository) DeleteEqualizer(id int) error {

	fnd := model.Equalizer{}
	if err := r.db.First(&fnd, id).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&fnd).Error; err != nil {
		return err
	}

	return nil
}

func (r Repository) EqualizerByName(name string, data interface{}) error {

	fnd := model.Equalizer{Name: name}

	if err := r.db.Where(&fnd).First(&data).Error; err != nil {
		return err
	}
	return nil
}
