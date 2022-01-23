package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type RequireItem struct {
	ItemId   int   `gorm:"type:int;not null" json:"type_id"`
	Quantity int64 `gorm:"type:int;not null" json:"quantity"`
}

type RequireItems []RequireItem

func (r *RequireItems) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("%v is not []byte", value)
	}

	return json.Unmarshal(str, &r)
}

func (r RequireItems) Value() (driver.Value, error) {
	str, err := json.Marshal(r)
	if err != nil {
		return nil, nil
	}

	return string(str), nil
}
