package model

import (
	"database/sql/driver"
	"encoding/json"
	"evelp/config/global"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

type BluePrint struct {
	BlueprintId int               `gorm:"type:int;not null;primary_key"`
	Products    ManufactProducts  `gorm:"type:text;not null"`
	Materials   ManufactMaterials `gorm:"type:text;not null"`
}

type ManufactProduct struct {
	ItemId   int
	Quantity int64
}

type ManufactMaterial struct {
	ItemId   int
	Quantity int64
}

type BluePrints []BluePrint

type ManufactProducts []ManufactProduct

type ManufactMaterials []ManufactMaterial

func (b BluePrints) Len() int { return len(b) }

func (b BluePrints) Less(i, j int) bool {
	return b[i].BlueprintId < b[j].BlueprintId
}

func (b BluePrints) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b *BluePrint) Empty() bool {
	if b.BlueprintId == 0 && b.Materials == nil && b.Products == nil {
		return true
	}
	return false
}

func (m *ManufactProducts) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return errors.Errorf("%v is not []byte", value)
	}

	return json.Unmarshal(str, &m)
}

func (m ManufactProducts) Value() (driver.Value, error) {
	str, err := json.Marshal(m)
	if err != nil {
		return nil, nil
	}

	return string(str), nil
}

func (m *ManufactMaterials) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return errors.Errorf("%v is not []byte", value)
	}

	return json.Unmarshal(str, &m)
}

func (m ManufactMaterials) Value() (driver.Value, error) {
	str, err := json.Marshal(m)
	if err != nil {
		return nil, nil
	}
	return string(str), nil
}

func GetBluePrint(blueprintItemId int) *BluePrint {
	var bluePrint BluePrint
	global.DB.Find(&bluePrint, blueprintItemId)
	return &bluePrint
}

func SaveBluePrint(bluePrint *BluePrint) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&bluePrint).Error; err != nil {
		return errors.Wrap(err, "save blueprint to DB failed")
	}
	return nil
}

func SaveBluePrints(bluePrints *BluePrints) error {
	for _, bluePrint := range *bluePrints {
		if err := SaveBluePrint(&bluePrint); err != nil {
			return err
		}
	}
	return nil
}
