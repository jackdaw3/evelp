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
		return errors.Wrapf(err, "redis marshal value %v error", value)
	}

	if err := global.Redis.Set(ctx, key, val, expirationTime).Err(); err != nil {
		return errors.Wrapf(err, "redis set key %v value %v error", key, val)
	}

	return nil
}

func Get(key string, dest interface{}) error {
	val, err := global.Redis.Get(ctx, key).Result()
	if err != nil {
		return errors.Wrapf(err, "redis get %v error", key)
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return errors.Wrapf(err, "redis unmarshal value %v error", val)
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
