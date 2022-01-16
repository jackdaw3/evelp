package global

import (
	"github.com/panjf2000/ants/v2"
	"gorm.io/gorm"
)

var (
	ENV       string
	WORKSPACE string
	DB        *gorm.DB
	ANTS      *ants.Pool
	LANGS     [6]string = [6]string{DE, EN, FR, JA, RU, ZH}
)
