package model

import (
	"evelp/config/global"

	"gorm.io/gorm/clause"
)

type Region struct {
	RegionId int  `gorm:"type:int;not null;primary_key;autoIncrement:false"`
	Name     Name `gorm:"type:text"`
}

type Regions []*Region

func (regions Regions) Len() int { return len(regions) }

func (regions Regions) Less(i, j int) bool { return regions[i].RegionId < regions[j].RegionId }

func (regions Regions) Swap(i, j int) { regions[i], regions[j] = regions[j], regions[i] }

func GetRegions() (*Regions, error) {
	var regions Regions
	result := global.DB.Find(&regions)
	return &regions, result.Error
}

func SaveRegion(region *Region) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&region).Error; err != nil {
		return err
	}
	return nil
}

func SaveRegions(regions *Regions) error {
	for _, region := range *regions {
		if err := SaveRegion(region); err != nil {
			return err
		}
	}
	return nil
}
