package service

import (
	"encoding/json"
	"evelp/model"
	"evelp/util/cache"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestOrders(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch(model.GetItem, func(int) (*model.Item, error) {
		return &model.Item{ItemId: 34, Name: model.Name{En: "Tritanium"}}, nil
	})

	monkey.Patch(model.GetStarSystem, func(int) (*model.StarSystem, error) {
		return &model.StarSystem{SystemId: 34, Name: model.Name{En: "Jita"}}, nil
	})

	mockOrders()

	orderService := NewOrderService(34, 1000002, false, 0.05)

	offers, err := orderService.Orders(true, "en")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(*offers))
	assert.Equal(t, "Tritanium", (*offers)[0].ItemName)

	offers, err = orderService.Orders(false, "en")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(*offers))

	orderService = NewOrderService(33, 1000002, true, 0.05)
	offers, err = orderService.Orders(true, "en")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(*offers))
	assert.Equal(t, "Tritanium", (*offers)[0].ItemName)
}

func TestHighestBuyPrice(t *testing.T) {
	defer monkey.UnpatchAll()
	mockOrders()

	orderService := NewOrderService(34, 1000002, false, 0.05)
	price, err := orderService.HighestBuyPrice()
	assert.NoError(t, err)
	assert.Equal(t, 4.14, price)
}

func TestLowestSellPrice(t *testing.T) {
	defer monkey.UnpatchAll()
	mockOrders()

	orderService := NewOrderService(34, 1000002, false, 0.05)
	price, err := orderService.LowestSellPrice()
	assert.NoError(t, err)
	assert.Equal(t, 5, int(price))
}

func mockOrders() {
	var (
		orders = &model.Orders{
			model.Order{
				OrderId:      6206535742,
				ItemId:       34,
				Duration:     90,
				SystemId:     30000142,
				Price:        5.4,
				VolumeRemain: 67762,
				VolumeTotal:  67762,
				IsBuyOrder:   false,
				LastUpdated:  time.Now(),
			},
			model.Order{
				OrderId:      6206524296,
				ItemId:       34,
				Duration:     3,
				SystemId:     30000142,
				Price:        4.45,
				VolumeRemain: 80948,
				VolumeTotal:  85510,
				IsBuyOrder:   false,
				LastUpdated:  time.Now(),
			},
			model.Order{
				OrderId:      6177179860,
				ItemId:       34,
				Duration:     90,
				SystemId:     30000142,
				Price:        5.04,
				VolumeRemain: 100000000,
				VolumeTotal:  100000000,
				IsBuyOrder:   false,
				LastUpdated:  time.Now(),
			},
			model.Order{
				OrderId:      6173986721,
				ItemId:       34,
				Duration:     90,
				SystemId:     30000142,
				Price:        3.2,
				VolumeRemain: 75000000,
				VolumeTotal:  75000000,
				IsBuyOrder:   true,
				LastUpdated:  time.Now(),
			},
			model.Order{
				OrderId:      6153010341,
				ItemId:       34,
				Duration:     90,
				SystemId:     30000142,
				Price:        4.14,
				VolumeRemain: 9915720,
				VolumeTotal:  10000000,
				IsBuyOrder:   true,
				LastUpdated:  time.Now(),
			},
		}

		bluePrint = &model.BluePrint{
			BlueprintId: 33,
			Products: model.ManufactProducts{
				model.ManufactProduct{
					ItemId:   34,
					Quantity: 1,
				},
			},
		}
	)

	monkey.Patch(cache.Get, func(key string, dest interface{}) error {
		val, err := json.Marshal(orders)
		if err != nil {
			return err
		}
		json.Unmarshal([]byte(val), dest)
		return nil
	})

	monkey.Patch(model.GetBluePrint, func(bluePrintId int) (*model.BluePrint, error) {
		if bluePrintId == bluePrint.BlueprintId {
			return bluePrint, nil
		}
		return nil, nil
	})
}
