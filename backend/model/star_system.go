package model

import (
	"evelp/config/global"
	"evelp/util/cache"
	"reflect"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

const (
	system_key        = "system"
	system_expiration = -1 * time.Second
)

type StarSystem struct {
	SystemId int  `gorm:"type:int;not null;primary_key;autoIncrement:false" json:"system_id"`
	Name     Name `gorm:"type:text" json:"name"`
}

type StarSystems []*StarSystem

func (s StarSystems) Len() int { return len(s) }

func (s StarSystems) Less(i, j int) bool {
	return s[i].SystemId < s[j].SystemId
}

func (s StarSystems) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func GetStarSystem(systemId int) (*StarSystem, error) {
	var starSystem StarSystem

	key := cache.Key(system_key, strconv.Itoa(systemId))
	exist := cache.Exist(key)

	if exist == nil {
		if err := cache.Get(key, &starSystem); err != nil {
			return nil, err
		}
		return &starSystem, nil
	} else {
		result := global.DB.First(&starSystem, systemId)
		if err := cache.Set(key, starSystem, system_expiration); err != nil {
			return nil, err
		}
		return &starSystem, result.Error
	}
}

func SaveStarSystem(starSystem *StarSystem) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&starSystem).Error; err != nil {
		return errors.Wrap(err, "save star system to DB failed")
	}

	key := cache.Key(system_key, strconv.Itoa(starSystem.SystemId))
	if err := cache.Set(key, *starSystem, system_expiration); err != nil {
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
