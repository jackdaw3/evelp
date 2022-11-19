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

const item_key = "item"

type Item struct {
	ItemId int     `gorm:"type:int;not null;primary_key;autoIncrement:false" json:"item_id"`
	Name   Name    `gorm:"type:text" yaml:"name" json:"name"`
	Volume float64 `gorm:"type:double" yaml:"volume" json:"volume"`
}

type Items []Item

func (items Items) Len() int { return len(items) }

func (items Items) Less(i, j int) bool { return items[i].ItemId < items[j].ItemId }

func (items Items) Swap(i, j int) { items[i], items[j] = items[j], items[i] }

func GetItem(id int) (*Item, error) {
	var item Item

	key := cache.Key(item_key, strconv.Itoa(id))
	if err := cache.Get(key, &item); err != nil {
		log.Debugf("failed to get item %d from cache: %s", id, err.Error())
		result := global.DB.First(&item, id)
		if err := cache.Set(key, item, global.Conf.Redis.ExpireTime.Model*time.Minute); err != nil {
			return nil, err
		}
		return &item, result.Error
	}

	return &item, nil
}

func SaveItem(item *Item) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&item).Error; err != nil {
		return errors.Wrapf(err, "failed to save item %d to DB", item.ItemId)
	}

	return nil
}

func SaveItems(items *Items) error {
	for _, item := range *items {
		if err := SaveItem(&item); err != nil {
			return err
		}
	}
	return nil
}
