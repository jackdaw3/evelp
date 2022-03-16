package service

import (
	"evelp/dto"
	"evelp/model"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestItemStatis(t *testing.T) {
	defer monkey.UnpatchAll()
	mockItemStatis()

	itemSerivce := NewItemStatisService(3414, 10000002, 0.05, "buy", 4.4, "en")
	statis, err := itemSerivce.ItemStatis(true)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(*statis))
	assert.Equal(t, 2, len((*statis)[0].Orderwrappers))
	assert.Equal(t, 530, (*statis)[0].AveUnitProfit)
	assert.Equal(t, "546 ~ 530", (*statis)[0].UnitProfitRange)

	statis, err = itemSerivce.ItemStatis(false)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(*statis))
	assert.Equal(t, 2, len((*statis)[0].Orderwrappers))
	assert.Equal(t, 758, (*statis)[0].AveUnitProfit)
	assert.Equal(t, "748 ~ 763", (*statis)[0].UnitProfitRange)
	assert.Equal(t, "1033 ~ 1033", (*statis)[1].UnitProfitRange)
}

func mockItemStatis() {
	monkey.Patch(model.GetOffer, func(offerId int) (*model.Offer, error) {
		return &model.Offer{OfferId: 3414, ItemId: 9943, Quantity: 1, IskCost: 5250000, LpCost: 5250}, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(&OrderService{}), "Orders", func(orderService *OrderService, isBuyOrder bool, lang string) (*dto.OrderDTOs, error) {
		if isBuyOrder {
			return &dto.OrderDTOs{
				dto.OrderDTO{OrderId: 6206571369, ItemId: 9943, VolumeRemain: 40, VolumeTotal: 40, Price: 8507000, IsBuyOrder: true},
				dto.OrderDTO{OrderId: 6200532940, ItemId: 9943, VolumeRemain: 1, VolumeTotal: 50, Price: 8600000, IsBuyOrder: true},
			}, nil
		} else {
			return &dto.OrderDTOs{
				dto.OrderDTO{OrderId: 6195025898, ItemId: 9943, VolumeRemain: 17, VolumeTotal: 22, Price: 9787000, IsBuyOrder: false},
				dto.OrderDTO{OrderId: 6192962730, ItemId: 9943, VolumeRemain: 7, VolumeTotal: 17, Price: 11270000, IsBuyOrder: false},
				dto.OrderDTO{OrderId: 6195122844, ItemId: 9943, VolumeRemain: 7, VolumeTotal: 12, Price: 9704000, IsBuyOrder: false},
			}, nil
		}
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(&OfferSerivce{}), "Offer", func(*OfferSerivce, int) (*dto.OfferDTO, error) {
		return &dto.OfferDTO{OfferId: 3414, ItemId: 9943, Quantity: 1, IskCost: 5250000, LpCost: 5250, MaterialCost: 100000}, nil
	})
}
