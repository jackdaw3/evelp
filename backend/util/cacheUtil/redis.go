package cacheUtil

import (
	"context"
	"encoding/json"
	"evelp/config/global"
	"strings"
	"time"
)

var ctx = context.Background()

func Set(key string, value interface{}, expirationTime time.Duration, groups ...string) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return global.REDIS.Set(ctx, groupKey(groups, key), val, expirationTime).Err()
}

func Get(key string, dest interface{}, groups ...string) error {
	val, err := global.REDIS.Get(ctx, groupKey(groups, key)).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func groupKey(groups []string, key string) string {
	var build strings.Builder

	for _, group := range groups {
		build.WriteString(group)
		build.WriteString(":")
	}
	build.WriteString(key)

	return build.String()
}
