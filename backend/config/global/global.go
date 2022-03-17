package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/panjf2000/ants/v2"
	"gorm.io/gorm"
)

var (
	Env       string
	Workspace string
	Conf      *Config
	DB        *gorm.DB
	Ants      *ants.Pool
	Redis     *redis.Client
	Langs     [6]string = [6]string{DE, EN, FR, JA, RU, ZH}
)
