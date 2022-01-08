package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Name struct {
	De string `yaml:"de"`
	En string `yaml:"en"`
	Fr string `yaml:"fr"`
	Ja string `yaml:"ja"`
	Ru string `yaml:"ru"`
	Zh string `yaml:"zh"`
}

func (name *Name) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}

	return json.Unmarshal(str, &name)
}

func (name Name) Value() (driver.Value, error) {
	str, err := json.Marshal(name)
	if err != nil {
		return nil, nil
	}

	return string(str), nil
}
