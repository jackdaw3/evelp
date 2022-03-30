package cache

import (
	"encoding/json"
	"evelp/config/global"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

const (
	requrestType   = "order"
	regionId       = "10000002"
	itemId         = "34"
	key            = "order:10000002:34"
	not_exist_key  = "order:10000002:35"
	expirationTime = time.Hour
)

type order_test struct {
	OrderId      int       `json:"order_id"`
	ItemId       int       `json:"type_id"`
	Issued       time.Time `json:"issued"`
	Duration     int       `json:"duration"`
	SystemId     int       `json:"system_id"`
	Price        float64   `json:"price"`
	VolumeRemain int64     `json:"volume_remain"`
	VolumeTotal  int64     `json:"volume_total"`
	IsBuyOrder   bool      `json:"is_buy_order"`
	LastUpdated  time.Time `json:"last_updated"`
}

var (
	server *miniredis.Miniredis
	order  *order_test
)

func setUp() error {

	issued, err := time.Parse(time.RFC3339, "2022-01-07T05:15:59Z")
	if err != nil {
		return err
	}

	order = &order_test{
		OrderId:      6173392220,
		ItemId:       28758,
		Issued:       issued,
		Duration:     90,
		SystemId:     30000142,
		Price:        26350000,
		VolumeRemain: 28,
		VolumeTotal:  30,
		IsBuyOrder:   true,
		LastUpdated:  time.Time{},
	}

	server, err = miniredis.Run()

	if err != nil {
		return err
	}

	global.Redis = redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	return nil
}

func TestGet(t *testing.T) {
	setUp()
	defer server.Close()

	b, err := json.Marshal(order)
	assert.NoError(t, err)
	server.Set(key, string(b))

	result := new(order_test)

	err = Get(key, result)
	assert.NoError(t, err)
	assert.Equal(t, order, result)

	err = Get(not_exist_key, result)
	assert.Equal(t, fmt.Sprintf("redis get %s error: redis: nil", not_exist_key), err.Error())
}

func TestSet(t *testing.T) {
	setUp()
	defer server.Close()

	b, err := json.Marshal(order)
	assert.NoError(t, err)

	err = Set(key, b, expirationTime)
	assert.NoError(t, err)

	assert.True(t, server.Exists(key))
}

func TestKey(t *testing.T) {
	val := Key(requrestType, regionId, itemId)
	assert.Equal(t, key, val)
}
