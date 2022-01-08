package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Description struct {
	De string `yaml:"de"`
	En string `yaml:"en"`
	Fr string `yaml:"fr"`
	Ja string `yaml:"ja"`
	Ru string `yaml:"ru"`
	Zh string `yaml:"zh"`
}

func (description *Description) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}

	return json.Unmarshal(str, &description)
}

func (description Description) Value() (driver.Value, error) {
	str, err := json.Marshal(description)
	if err != nil {
		return nil, nil
	}

	return string(str), nil
}
