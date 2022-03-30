package initial

import (
	"context"
	"evelp/config/global"
	"evelp/log"
	"evelp/util/crypto"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

var ctx = context.Background()

func initRedis() error {
	password, err := crypto.Decrypt(global.Conf.Redis.Password, global.Conf.Crypto.KeyPath)
	if err != nil {
		return errors.WithMessage(err, "decode redis password error")
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

	global.Redis = rdb
	return nil
}
