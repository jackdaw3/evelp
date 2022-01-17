package cache

import (
	"context"
	"encoding/json"
	"evelp/config/global"
	"strings"
	"time"
)

var ctx = context.Background()

func Set(key string, value interface{}, expirationTime time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return global.REDIS.Set(ctx, key, val, expirationTime).Err()
}

func Get(key string, dest interface{}, groups ...string) error {
	val, err := global.REDIS.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
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
