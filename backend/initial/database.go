package initial

import (
	"encoding/base64"
	"evelp/config/global"
	"evelp/model"
	"fmt"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func database() error {
	b, err := base64.StdEncoding.DecodeString(global.Conf.MySQL.Password)
	if err != nil {
		return fmt.Errorf("decode database password failed: %v", err)
	}
	password := string(b)

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
	log.Infof("connect to mysql: %s", dsn)

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
			return err
		}
	}

	return nil
}
