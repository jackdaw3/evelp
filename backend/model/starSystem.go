package model

import (
	"evelp/config/global"

	"gorm.io/gorm/clause"
)

type StarSystem struct {
	SystemId int  `gorm:"type:int;not null;primary_key;autoIncrement:false"`
	Name     Name `gorm:"type:text"`
}

type StarSystems []*StarSystem

func (starSystems StarSystems) Len() int { return len(starSystems) }

func (starSystems StarSystems) Less(i, j int) bool {
	return starSystems[i].SystemId < starSystems[j].SystemId
}

func (starSystems StarSystems) Swap(i, j int) {
	starSystems[i], starSystems[j] = starSystems[j], starSystems[i]
}

func GetStarSystem(systemId int) (*StarSystem, error) {
	var starSystem StarSystem
	result := global.DB.First(&starSystem, systemId)
	return &starSystem, result.Error
}

func SaveStarSystem(starSystem *StarSystem) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&starSystem).Error; err != nil {
		return err
	}
	return nil
}

func SaveStarSystems(starSystems *StarSystems) error {
	for _, starSystem := range *starSystems {
		if err := SaveStarSystem(starSystem); err != nil {
			return err
		}
	}
	return nil
}

func IsStarSystemExist(systemId int) (bool, error) {
	var starSystem StarSystem

	count := int64(0)
	err := global.DB.Model(&starSystem).Where("system_id = ?", systemId).Count(&count).Error
	if err != nil {
		return false, err
	}

	exists := count > 0
	return exists, nil
}
