package initial

import (
	"context"
	"encoding/base64"
	"evelp/config/global"
	"fmt"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

func initRedis() error {
	b, err := base64.StdEncoding.DecodeString(global.Conf.Redis.Password)
	if err != nil {
		return fmt.Errorf("decode redis password failed: %v", err)
	}
	password := string(b)

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
