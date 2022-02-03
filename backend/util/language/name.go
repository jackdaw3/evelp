package language

import (
	"evelp/config/global"
	"evelp/model"
)

func Name(lang string, Name model.Name) string {
	var name string

	switch lang {
	case global.DE:
		name = Name.De
	case global.EN:
		name = Name.En
	case global.FR:
		name = Name.Fr
	case global.JA:
		name = Name.Ja
	case global.RU:
		name = Name.Ru
	case global.ZH:
		name = Name.Zh
	}

	return name
}
