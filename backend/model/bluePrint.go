package model

import (
	"database/sql/driver"
	"encoding/json"
	"evelp/configs/global"
	"fmt"

	"gorm.io/gorm/clause"
)

type BluePrint struct {
	BlueprintId int               `gorm:"type:int;not null;primary_key"`
	Products    ManufactProducts  `gorm:"type:text;not null"`
	Materials   ManufactMaterials `gorm:"type:text;not null"`
}

type ManufactProduct struct {
	ItemId   int
	Quantity int
}

type ManufactMaterial struct {
	ItemId   int
	Quantity int
}

type BluePrints []BluePrint

type ManufactProducts []ManufactProduct

type ManufactMaterials []ManufactMaterial

func (bluePrints BluePrints) Len() int { return len(bluePrints) }

func (bluePrints BluePrints) Less(i, j int) bool {
	return bluePrints[i].BlueprintId < bluePrints[j].BlueprintId
}

func (bluePrints BluePrints) Swap(i, j int) {
	bluePrints[i], bluePrints[j] = bluePrints[j], bluePrints[i]
}

func (b *BluePrint) Empty() bool {
	if b.BlueprintId == 0 && b.Materials == nil && b.Products == nil {
		return true
	}
	return false
}

func (manufactProducts *ManufactProducts) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}

	return json.Unmarshal(str, &manufactProducts)
}

func (manufactProducts ManufactProducts) Value() (driver.Value, error) {
	str, err := json.Marshal(manufactProducts)
	if err != nil {
		return nil, nil
	}

	return string(str), nil
}

func (manufactMaterials *ManufactMaterials) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}

	return json.Unmarshal(str, &manufactMaterials)
}

func (manufactMaterials ManufactMaterials) Value() (driver.Value, error) {
	str, err := json.Marshal(manufactMaterials)
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
		return err
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
