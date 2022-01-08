package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type RequireItem struct {
	ItemId   int `gorm:"type:int;not null" json:"type_id"`
	Quantity int `gorm:"type:int;not null" json:"quantity"`
}

type RequireItems []RequireItem

func (requireItems *RequireItems) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}

	return json.Unmarshal(str, &requireItems)
}

func (requireItems RequireItems) Value() (driver.Value, error) {
	str, err := json.Marshal(requireItems)
	if err != nil {
		return nil, nil
	}

	return string(str), nil
}
