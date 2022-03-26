package model

import (
	"evelp/config/global"
	"evelp/util/cache"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

const (
	item_key        = "item"
	item_expiration = -1 * time.Second
)

//easyjson:json
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
	exist := cache.Exist(key)

	if exist == nil {
		if err := cache.Get(key, &item); err != nil {
			return nil, err
		}
		return &item, nil
	} else {
		result := global.DB.First(&item, id)
		if err := cache.Set(key, &item, item_expiration); err != nil {
			return nil, err
		}
		return &item, result.Error
	}
}

func SaveItem(item *Item) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&item).Error; err != nil {
		return errors.Wrap(err, "save item to DB failed")
	}

	key := cache.Key(item_key, strconv.Itoa(item.ItemId))
	if err := cache.Set(key, item, item_expiration); err != nil {
		return err
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
