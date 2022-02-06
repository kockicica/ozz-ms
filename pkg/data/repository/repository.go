package repository

import (
	"fmt"
	"strings"

	"ozz-ms/pkg/data/model"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Repository struct {
	db *gorm.DB
}

type RepositoryConfig struct {
	Dsn     string
	Verbose bool
	Logger  gormlogger.Interface
}

func NewRepository(cfg RepositoryConfig) (*Repository, error) {

	var err error
	var db *gorm.DB
	var dialect gorm.Dialector

	dsnParts := strings.Split(cfg.Dsn, ":")
	if len(dsnParts) == 2 {
		// db + dsn
		dialect, err = createDialector(dsnParts[0], dsnParts[1])
		if err != nil {
			return nil, err
		}
	} else {
		dialect, err = createDialector("sqlite", cfg.Dsn)
		if err != nil {
			return nil, err
		}
	}

	repo := new(Repository)

	gormCfg := gorm.Config{}
	if cfg.Logger != nil {
		gormCfg.Logger = cfg.Logger
	} else {
		gormCfg.Logger = gormlogger.Default
	}

	if cfg.Verbose {
		gormCfg.Logger = gormCfg.Logger.LogMode(gormlogger.Info)
	}

	db, err = gorm.Open(dialect, &gormCfg)

	//db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}

	models := []interface{}{
		&model.Shift{},
		&model.Category{},
		&model.AudioRecording{},
		&model.Disposition{},
		&model.DispositionPlayed{},
		&model.User{},
		&model.Schedule{},
	}

	if err = db.AutoMigrate(models...); err != nil {
		return nil, err
	}

	if err = initCategories(db); err != nil {
		return nil, err
	}

	if err = initShifts(db); err != nil {
		return nil, err
	}

	if err = initUsers(db); err != nil {
		return nil, err
	}

	repo.db = db

	return repo, nil

}

func initCategories(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Category{}).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return nil
	}
	var predefinedCategories = []model.Category{
		{
			Name:  "REKLAME",
			Path:  "reklame",
			Order: 1,
		},
		{
			Name:  "ŠPICE",
			Path:  "spice",
			Order: 2,
		},
		{
			Name:  "UPADICE",
			Path:  "upadice",
			Order: 3,
		},
		{
			Name:  "MASKE",
			Path:  "maske",
			Order: 3,
		},
		{
			Name:  "DŽINGLOVI",
			Path:  "dzinglovi",
			Order: 4,
		},
		{
			Name:  "SLOBODNA",
			Path:  "slobodna",
			Order: 5,
		},
	}

	for _, c := range predefinedCategories {
		if err := db.Create(&c).Error; err != nil {
			return err
		}
	}
	return nil
}

func initShifts(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Shift{}).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return nil
	}

	predefinedShifts := []model.Shift{
		{Name: "Smena I", Active: true, Order: 2},
		{Name: "Smena II", Active: true, Order: 3},
		{Name: "Smena III", Active: true, Order: 4},
		{Name: "Smena IV", Active: true, Order: 1},
	}

	for _, s := range predefinedShifts {
		if err := db.Create(&s).Error; err != nil {
			return err
		}
	}

	return nil
}

func initUsers(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return nil
	}

	predefinedUsers := []model.User{
		{Username: "maki", Password: "maki", Level: model.Admin},
		{Username: "taki", Password: "taki", Level: model.Regular},
		{Username: "laki", Password: "laki", Level: model.Regular},
		{Username: "caki", Password: "caki", Level: model.Regular},
	}

	for _, u := range predefinedUsers {
		if err := u.SetPassword(u.Username); err != nil {
			return err
		}
		if err := db.Create(&u).Error; err != nil {
			return err
		}
	}

	return nil
}

func createDialector(db, dsn string) (gorm.Dialector, error) {
	switch db {
	case "sqlite":
		return sqlite.Open(dsn), nil
	case "mysql":
		return mysql.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unable to find gorm dialector for db: %s", db)
	}
}
