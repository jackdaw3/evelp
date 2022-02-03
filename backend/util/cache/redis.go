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
		return errors.Wrap(err, "redis set marshal failed")
	}
	return global.REDIS.Set(ctx, key, val, expirationTime).Err()
}

func Get(key string, dest interface{}) error {
	val, err := global.REDIS.Get(ctx, key).Result()
	if err != nil {
		return errors.Wrap(err, "redis get result failed")
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return errors.Wrap(err, "redis get unmarshal failed")
	}
	return nil
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
