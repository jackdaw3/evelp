package service

import (
	"evelp/model"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

var (
	offer = model.Offer{
		OfferId:  3439,
		ItemId:   17703,
		Quantity: 1,
		IskCost:  0,
		LpCost:   10000,
		AkCost:   0,
		RequireItems: model.RequireItems{
			model.RequireItem{ItemId: 17795, Quantity: 1},
			model.RequireItem{ItemId: 591, Quantity: 1},
		},
		CorporationIds: model.CorporationIds{1000019, 1000064},
		IsBluePrint:    false,
	}

	bluePrintOffer = model.Offer{
		OfferId:        4335,
		ItemId:         17874,
		Quantity:       1,
		IskCost:        10000000,
		LpCost:         30000,
		AkCost:         0,
		CorporationIds: model.CorporationIds{1000019, 1000064},
		IsBluePrint:    true,
	}

	bluePrint = model.BluePrint{
		BlueprintId: 17874,
		Products: model.ManufactProducts{
			model.ManufactProduct{ItemId: 17873, Quantity: 1},
		},
		Materials: model.ManufactMaterials{
			model.ManufactMaterial{ItemId: 34, Quantity: 1111},
		},
	}
)

func TestOffer(t *testing.T) {
	defer monkey.UnpatchAll()
	mockOffers()

	monkey.Patch(model.GetOffer, func(offerId int) (*model.Offer, error) {
		return &offer, nil
	})

	offerService := NewOfferSerivce(1000002, 0.05, 7, "buy", "sell", 0, "en")
	offerDTO, err := offerService.Offer(offer.ItemId)

	assert.NoError(t, err)
	assert.Equal(t, 17703, offerDTO.ItemId)
}

func TestOffers(t *testing.T) {
	defer monkey.UnpatchAll()
	offerService := NewOfferSerivce(1000002, 0.05, 7, "buy", "sell", 0, "en")
	mockOffers()

	offers, err := offerService.Offers(1000019)

	expectedUnitProfit1 := (64000000 - 10000000 - 4.45*1111) / 30000
	expectedUnitProfit2 := (12300000 - (120600 + 400000)) / 10000

	assert.NoError(t, err)
	assert.Equal(t, 2, len(*offers))
	assert.Equal(t, "System Scanner I Blueprint", (*offers)[0].ItemName)
	assert.Equal(t, "Imperial Navy Slicer", (*offers)[1].ItemName)
	assert.Equal(t, int(expectedUnitProfit1), (*offers)[0].UnitProfit)
	assert.Equal(t, int(expectedUnitProfit2), (*offers)[1].UnitProfit)
}

func mockOffers() {
	monkey.Patch(model.GetOffersByCorporation, func(int) (*model.Offers, error) {
		return &model.Offers{&offer, &bluePrintOffer}, nil
	})

	monkey.Patch(model.GetBluePrint, func(int) (*model.BluePrint, error) {
		return &bluePrint, nil
	})

	monkey.Patch(model.GetItem, func(id int) (*model.Item, error) {
		switch id {
		case 34:
			return &model.Item{ItemId: 34, Name: model.Name{En: "Tritanium"}}, nil
		case 591:
			return &model.Item{ItemId: 591, Name: model.Name{En: "Tormentor"}}, nil
		case 17703:
			return &model.Item{ItemId: 17703, Name: model.Name{En: "Imperial Navy Slicer"}}, nil
		case 17795:
			return &model.Item{ItemId: 17795, Name: model.Name{En: "Amarr MIY-1 Nexus Chip"}}, nil
		case 17873:
			return &model.Item{ItemId: 17873, Name: model.Name{En: "System Scanner I"}}, nil
		case 17874:
			return &model.Item{ItemId: 17874, Name: model.Name{En: "System Scanner I Blueprint"}}, nil
		}
		return nil, nil
	})

	monkey.Patch(NewOrderService, func(itemId int, regionId int, isBluePrint bool, scope float64) *OrderService {
		switch itemId {
		case 34:
			return &OrderService{34, 1000002, false, 0.05}
		case 591:
			return &OrderService{591, 1000002, false, 0.05}
		case 17703:
			return &OrderService{17703, 1000002, false, 0.05}
		case 17795:
			return &OrderService{17795, 1000002, false, 0.05}
		case 17873:
			return &OrderService{17873, 1000002, false, 0.05}
		case 17874:
			return &OrderService{17874, 1000002, true, 0.05}
		}
		return nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(&OrderService{}), "LowestSellPrice", func(orderService *OrderService) (float64, error) {
		switch orderService.itemId {
		case 34:
			return 4.45, nil
		case 591:
			return 400000, nil
		case 17703:
			return 13300000, nil
		case 17795:
			return 120600, nil
		case 17873:
			return 80000000, nil
		case 17874:
			return 80000000, nil
		}
		return 0, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(&OrderService{}), "HighestBuyPrice", func(orderService *OrderService) (float64, error) {
		switch orderService.itemId {
		case 34:
			return 4.25, nil
		case 591:
			return 300000, nil
		case 17703:
			return 12300000, nil
		case 17795:
			return 110600, nil
		case 17873:
			return 64000000, nil
		case 17874:
			return 64000000, nil
		}
		return 0, nil
	})

	monkey.Patch(NewItemHistoryService, func(itemId int, regionId int, isBluePrint bool) *ItemHistoryService {
		switch itemId {
		case 34:
			return &ItemHistoryService{34, 1000002, false}
		case 591:
			return &ItemHistoryService{591, 1000002, false}
		case 17703:
			return &ItemHistoryService{17703, 1000002, false}
		case 17795:
			return &ItemHistoryService{17795, 1000002, false}
		case 17873:
			return &ItemHistoryService{17873, 1000002, false}
		case 17874:
			return &ItemHistoryService{17874, 1000002, true}
		}
		return nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(&ItemHistoryService{}), "AverageVolume", func(*ItemHistoryService, int) (int64, error) {
		return 10, nil
	})
}
