package cache

import (
	"encoding/json"
	"evelp/config/global"
	"evelp/model"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

const (
	requrestType   = "order"
	regionId       = "10000002"
	itemId         = "34"
	key            = "order:10000002:34"
	expirationTime = time.Hour
)

var order *model.Order

func setUp() error {
	issued, err := time.Parse(time.RFC3339, "2022-01-07T05:15:59Z")
	if err != nil {
		return err
	}

	order = &model.Order{
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

	return nil
}

func TestGet(t *testing.T) {
	redis, mock := redismock.NewClientMock()
	global.REDIS = redis
	setUp()

	b, err := json.Marshal(order)
	assert.NoError(t, err)
	mock.ExpectGet(key).SetVal(string(b))

	result := new(model.Order)
	err = Get(key, result)
	assert.NoError(t, err)

	assert.Equal(t, order, result)
}

func TestSet(t *testing.T) {
	redis, mock := redismock.NewClientMock()
	global.REDIS = redis

	setUp()
	b, err := json.Marshal(order)
	assert.NoError(t, err)
	mock.ExpectSet(key, b, expirationTime).SetVal("OK")

	err = Set(key, order, expirationTime)
	assert.NoError(t, err)
}

func TestKey(t *testing.T) {
	val := Key(requrestType, regionId, itemId)
	assert.Equal(t, key, val)
}
