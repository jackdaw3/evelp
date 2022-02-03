package model

import (
	"evelp/config/global"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

type Item struct {
	ItemId int     `gorm:"type:int;not null;primary_key;autoIncrement:false"`
	Name   Name    `gorm:"type:text" yaml:"name"`
	Volume float64 `gorm:"type:double" yaml:"volume"`
}

type Items []Item

func (items Items) Len() int { return len(items) }

func (items Items) Less(i, j int) bool { return items[i].ItemId < items[j].ItemId }

func (items Items) Swap(i, j int) { items[i], items[j] = items[j], items[i] }

func GetItem(id int) (*Item, error) {
	var item Item
	result := global.DB.First(&item, id)
	return &item, result.Error
}

func SaveItem(item *Item) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&item).Error; err != nil {
		return errors.Wrap(err, "save item to DB failed")
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
