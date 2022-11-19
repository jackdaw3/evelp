package model

import (
	"evelp/config/global"
	"evelp/log"
	"evelp/util/cache"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

const (
	corporation_key         = "corporation"
	corporations_key        = "corporation:corporations"
	corporation_faction_key = "corporation:faction"
)

type Corporation struct {
	CorporationId int  `gorm:"type:int;not null;primary_key;autoIncrement:false"`
	FactionId     int  `gorm:"type:int" yaml:"factionID"`
	Name          Name `gorm:"type:text" yaml:"nameID"`
}

type Corporations []Corporation

func (c Corporations) Len() int { return len(c) }

func (c Corporations) Less(i, j int) bool {
	return c[i].CorporationId < c[j].CorporationId
}

func (c Corporations) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

func GetCorporation(corporationId int) (*Corporation, error) {
	var corporation Corporation

	key := cache.Key(corporation_key, strconv.Itoa(corporationId))
	if err := cache.Get(key, &corporation); err != nil {
		log.Debugf("failed to get corporation %d from cache: %s", corporationId, err.Error())
		result := global.DB.First(&corporation, corporationId)
		if err := cache.Set(key, corporation, global.Conf.Redis.ExpireTime.Model*time.Minute); err != nil {
			return nil, err
		}
		return &corporation, result.Error
	}

	return &corporation, nil
}

func GetCorporations() (*Corporations, error) {
	var corporations Corporations

	if err := cache.Get(corporations_key, &corporations); err != nil {
		log.Debugf("failed to get all corporations from cache: %s", err.Error())
		result := global.DB.Find(&corporations)
		if err := cache.Set(corporations_key, corporations, global.Conf.Redis.ExpireTime.Model*time.Minute); err != nil {
			return nil, err
		}
		return &corporations, result.Error
	}

	return &corporations, nil
}

func GetCorporationsByFaction(factionId int) (*Corporations, error) {
	var corporations Corporations

	key := cache.Key(corporation_faction_key, strconv.Itoa(factionId))
	if err := cache.Get(key, &corporations); err != nil {
		log.Debugf("failed to get faction %d's corporations from cache: %s", factionId, err.Error())
		result := global.DB.Where("faction_id = ?", factionId).Find(&corporations)
		if err := cache.Set(key, &corporations, global.Conf.Redis.ExpireTime.Model*time.Minute); err != nil {
			return nil, err
		}
		return &corporations, result.Error
	}

	return &corporations, nil
}

func SaveCorporation(corporation *Corporation) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&corporation).Error; err != nil {
		return errors.Wrap(err, "failed to save corporation to DB ")
	}

	return nil
}

func SaveCorporations(corporations *Corporations) error {
	for _, corporation := range *corporations {
		if err := SaveCorporation(&corporation); err != nil {
			return err
		}
	}
	return nil
}
