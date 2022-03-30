package model

import (
	"evelp/config/global"
	"reflect"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

type Region struct {
	RegionId int  `gorm:"type:int;not null;primary_key;autoIncrement:false"`
	Name     Name `gorm:"type:text"`
}

type Regions []*Region

func (r Regions) Len() int { return len(r) }

func (r Regions) Less(i, j int) bool { return r[i].RegionId < r[j].RegionId }

func (r Regions) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func GetRegion(regionId int) (*Region, error) {
	var region Region
	result := global.DB.First(&region, regionId)
	return &region, result.Error
}

func GetRegions() (*Regions, error) {
	var regions Regions
	result := global.DB.Find(&regions)
	return &regions, result.Error
}

func SaveRegion(region *Region) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&region).Error; err != nil {
		return errors.Wrapf(err, "failed to save region %d to DB", region.RegionId)
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

func (r Region) IsExist() (bool, error) {
	var region Region

	count := int64(0)
	err := global.DB.Model(&region).Where("region_id = ?", r.RegionId).Count(&count).Error
	if err != nil {
		return false, err
	}

	exists := count > 0
	return exists, nil
}

func (r Region) IsVaild() (bool, error) {
	region, err := GetRegion(r.RegionId)
	if err != nil {
		return false, err
	}

	value := reflect.ValueOf(region.Name)
	langsCount := value.NumField()
	for i := 0; i < langsCount; i++ {
		field := value.Field(i)
		if len(field.String()) == 0 {
			return false, nil
		}
	}

	return true, nil
}
