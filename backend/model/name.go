package model

import (
	"database/sql/driver"
	"encoding/json"
	"evelp/config/global"

	"github.com/pkg/errors"
)

type Name struct {
	De string `yaml:"de" json:"de"`
	En string `yaml:"en" json:"en"`
	Fr string `yaml:"fr" json:"fr"`
	Ja string `yaml:"ja" json:"ja"`
	Ru string `yaml:"ru" json:"ru"`
	Zh string `yaml:"zh" json:"zh"`
}

func (n *Name) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return errors.Errorf("%v is not []byte", value)
	}

	return json.Unmarshal(str, &n)
}

func (n Name) Value() (driver.Value, error) {
	str, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}

	return string(str), nil
}

func (n *Name) Lang(lang string) string {
	var val string

	switch lang {
	case global.DE:
		val = n.De
	case global.EN:
		val = n.En
	case global.FR:
		val = n.Fr
	case global.JA:
		val = n.Ja
	case global.RU:
		val = n.Ru
	case global.ZH:
		val = n.Zh
	}

	return val
}
