package initial

import (
	"context"
	"evelp/config/global"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

func initRedis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     global.Conf.Redis.Address,
		Password: global.Conf.Redis.Password,
		DB:       global.Conf.Redis.Database,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	log.Infof("%v! Connect to redis:%s", pong, global.Conf.Redis.Address)

	global.REDIS = rdb
	return nil
}
