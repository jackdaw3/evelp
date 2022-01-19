package initial

import (
	"context"
	"evelp/config/global"
	"evelp/util/crypto"
	"fmt"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

func initRedis() error {
	password, err := crypto.Decrypt(global.Conf.MySQL.Password, global.Conf.Crypto.KeyPath)
	if err != nil {
		return fmt.Errorf("decode database password failed: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     global.Conf.Redis.Address,
		Password: password,
		DB:       global.Conf.Redis.Database,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	log.Infof("%v connect to redis server: %s", pong, global.Conf.Redis.Address)

	global.REDIS = rdb
	return nil
}
