package model

import (
	"evelp/config/global"

	"gorm.io/gorm/clause"
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

func GetCorporation(id int) (*Corporation, error) {
	var corporation Corporation
	result := global.DB.First(&corporation, id)
	return &corporation, result.Error
}

func GetCorporations() (*Corporations, error) {
	var corporations Corporations
	result := global.DB.Find(&corporations)
	return &corporations, result.Error
}

func GetCorporationsByFaction(factionId int) (*Corporations, error) {
	var corporations Corporations
	result := global.DB.Where("faction_id <> ?", factionId).Find(&corporations)
	return &corporations, result.Error
}

func SaveCorporation(corporation *Corporation) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&corporation).Error; err != nil {
		return err
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
