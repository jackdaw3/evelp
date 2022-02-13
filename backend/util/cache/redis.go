package cache

import (
	"context"
	"encoding/json"
	"evelp/config/global"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var ctx = context.Background()

func Set(key string, value interface{}, expirationTime time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return errors.Wrapf(err, "redis marshal value %v failed", value)
	}

	if err := global.REDIS.Set(ctx, key, val, expirationTime).Err(); err != nil {
		return errors.Wrapf(err, "redis set key %v value %v failed", key, val)
	}

	return nil
}

func Get(key string, dest interface{}) error {
	if err := Exist(key); err != nil {
		return err
	}

	val, err := global.REDIS.Get(ctx, key).Result()
	if err != nil {
		return errors.Wrapf(err, "redis get %v failed", key)
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return errors.Wrapf(err, "redis unmarshal value %v failed", val)
	}
	return nil
}

func Exist(key string) error {
	val, err := global.REDIS.Exists(ctx, key).Result()
	if err != nil {
		return errors.Wrapf(err, "redis check %v exist failed", key)
	}

	if val == 1 {
		return nil
	} else {
		return errors.Errorf("redis key %v not exist", key)
	}
}

func Key(args ...string) string {
	var build strings.Builder
	for index, arg := range args {
		build.WriteString(arg)
		if !(index == len(args)-1) {
			build.WriteString(":")
		}
	}
	return build.String()
}
