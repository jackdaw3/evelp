package global

import "gorm.io/gorm"

var (
	ENV       string
	WORKSPACE string
	DB        *gorm.DB
)
