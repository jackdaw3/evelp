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

func (n *Name) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("%v is not []byte", value)
	}

	return json.Unmarshal(str, &n)
}

func (n Name) Value() (driver.Value, error) {
	str, err := json.Marshal(n)
	if err != nil {
		return nil, nil
	}

	return string(str), nil
}
