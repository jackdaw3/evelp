package model

import (
	"evelp/config/global"
	"reflect"

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

func (s StarSystem) IsExist() (bool, error) {
	var starSystem StarSystem

	count := int64(0)
	err := global.DB.Model(&starSystem).Where("system_id = ?", s.SystemId).Count(&count).Error
	if err != nil {
		return false, err
	}

	exists := count > 0
	return exists, nil
}

func (s StarSystem) IsVaild() (bool, error) {
	starSystem, err := GetStarSystem(s.SystemId)
	if err != nil {
		return false, err
	}

	value := reflect.ValueOf(starSystem.Name)
	langsCount := value.NumField()
	for i := 0; i < langsCount; i++ {
		field := value.Field(i)
		if len(field.String()) == 0 {
			return false, nil
		}
	}

	return true, nil
}
