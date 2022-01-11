package initial

import (
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		global.Conf.MySQL.UserName,
		global.Conf.MySQL.Password,
		global.Conf.MySQL.Host,
		global.Conf.MySQL.Port,
		global.Conf.MySQL.Database,
		global.Conf.MySQL.Charset,
		url.QueryEscape(global.Conf.MySQL.Loc))

	log.Info("Connect to mysql:", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(global.Conf.MySQL.MaxIdleConn)
	sqlDB.SetMaxOpenConns(global.Conf.MySQL.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(global.Conf.MySQL.ConnMaxLifeTime))

	global.DB = db

	return nil
}

func autoMigrate() error {
	log.Info("Auto migrate db tables.")

	models := []interface{}{&model.Item{}, &model.Faction{}, &model.Corporation{}, &model.Offer{}, &model.BluePrint{}, &model.Region{}, &model.StarSystem{}}

	for _, m := range models {
		if err := global.DB.AutoMigrate(m); err != nil {
			return err
		}
	}

	return nil
}
