package initial

import (
	"evelp/config/global"
	"evelp/log"
	"evelp/model"
	"evelp/util/crypto"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func database() error {
	password, err := crypto.Decrypt(global.Conf.MySQL.Password, global.Conf.Crypto.KeyPath)
	if err != nil {
		return errors.WithMessage(err, "decode database password error")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		global.Conf.MySQL.UserName,
		password,
		global.Conf.MySQL.Host,
		global.Conf.MySQL.Port,
		global.Conf.MySQL.Database,
		global.Conf.MySQL.Charset,
		url.QueryEscape(global.Conf.MySQL.Loc))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	log.Infof("connect to mysql: %s", fmt.Sprintf("%s:%s", global.Conf.MySQL.Host, global.Conf.MySQL.Port))

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(global.Conf.MySQL.MaxIdleConn)
	sqlDB.SetMaxOpenConns(global.Conf.MySQL.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(global.Conf.MySQL.ConnMaxLifeTime))

	global.DB = db

	if global.Conf.MySQL.AutoMigrate {
		if err := autoMigrate(); err != nil {
			return err
		}
	}

	return nil
}

func autoMigrate() error {
	log.Info("auto migrate DB tables")

	models := []interface{}{&model.Item{}, &model.Faction{}, &model.Corporation{}, &model.Offer{}, &model.BluePrint{}, &model.Region{}, &model.StarSystem{}}

	for _, m := range models {
		if err := global.DB.AutoMigrate(m); err != nil {
			return errors.Wrap(err, "auto migrate db tables error")
		}
	}

	return nil
}
