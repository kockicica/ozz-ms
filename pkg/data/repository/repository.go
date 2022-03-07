package repository

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

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

	u, err := url.Parse(cfg.Dsn)
	if err != nil {
		return nil, err
	}

	dialect, err := createDialector(u)

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
		//&model.Disposition{},
		//&model.DispositionPlayed{},
		&model.User{},
		&model.Schedule{},
		&model.Equalizer{},
		&model.EmitLog{},
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

func createDialector(dbUrl *url.URL) (gorm.Dialector, error) {

	var err error

	switch dbUrl.Scheme {
	case "mysql":
		if dbUrl.User == nil {
			return nil, fmt.Errorf("there is no mssql authentication data")
		}
		pass, passOk := dbUrl.User.Password()
		if !passOk {
			return nil, fmt.Errorf("there is no mssql password supplied")
		}

		cn := fmt.Sprintf("%s:%s@tcp(%s)%s", dbUrl.User.Username(), pass, dbUrl.Host, dbUrl.Path)
		if dbUrl.RawQuery != "" {
			cn += fmt.Sprintf("?%s", dbUrl.RawQuery)
		}
		return mysql.Open(cn), nil
	case "sqlite":
		fallthrough
	default:
		var absPath string
		lpath := dbUrl.Path
		if lpath[0] == '/' {
			lpath = lpath[1:]
		}
		if !filepath.IsAbs(lpath) {
			cn := filepath.Join(dbUrl.Host, lpath)
			absPath, err = filepath.Abs(cn)
			if err != nil {
				return nil, err
			}
		} else {
			absPath = lpath
		}
		dir, _ := filepath.Split(absPath)
		if err = os.MkdirAll(dir, os.ModeDir); err != nil {
			return nil, err
		}
		absPath += "?_pragma=foreign_keys(1)"
		if dbUrl.RawQuery != "" {
			absPath += fmt.Sprintf("&%s", dbUrl.RawQuery)
		}
		return sqlite.Open(absPath), nil
	}
}
