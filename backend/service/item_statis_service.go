package service

import (
	"evelp/dto"
	"evelp/model"
	"fmt"
	"sort"
)

type ItemStatisService struct {
	offerId       int
	regionId      int
	scope         float64
	materialPrice string
	lang          string
}

func NewItemStatisService(offerId, regionId int, scope float64, materialPrice string, lang string) *ItemStatisService {
	return &ItemStatisService{offerId, regionId, scope, materialPrice, lang}
}

func (is *ItemStatisService) ItemStatis(isBuyOrder bool) (*dto.ItemStatisDTOs, error) {
	offer, err := model.GetOffer(is.offerId)
	if err != nil {
		return nil, err
	}

	var itemId int
	if offer.IsBluePrint {
		bluePrint := model.GetBluePrint(offer.ItemId)
		product, err := model.GetItem(bluePrint.Products[0].ItemId)
		if err != nil {
			return nil, err
		}
		itemId = product.ItemId
	} else {
		itemId = offer.ItemId
	}
	orderService := NewOrderService(itemId, is.regionId, is.scope)
	orders, err := orderService.Orders(isBuyOrder, is.lang)
	if err != nil {
		return nil, err
	}

	offerService := NewOfferSerivce(is.regionId, is.scope, 0, "", is.materialPrice, is.lang)
	offerDTO, err := offerService.Offer(is.offerId)
	if err != nil {
		return nil, err
	}

	unitCost := (offerDTO.MaterialCost + offerDTO.IskCost) / float64(offerDTO.Quantity)
	unitLpCost := offerDTO.LpCost / offerDTO.Quantity

	var orderwList dto.OrderDTOWrappers
	for _, order := range *orders {
		orderw := new(dto.OrderDTOWrapper)
		orderw.OrderDTO = order
		orderw.Income = order.Price * float64(order.VolumeRemain)
		orderw.Cost = unitCost * float64(order.VolumeRemain)
		orderw.Profit = orderw.Income - orderw.Cost
		orderw.UnitProfit = int(orderw.Profit / float64(int64(unitLpCost)*order.VolumeRemain))
		orderwList = append(orderwList, orderw)
	}
	if isBuyOrder {
		sort.Sort(sort.Reverse(orderwList))
	} else {
		sort.Sort(orderwList)
	}

	var itemStatisDTOS dto.ItemStatisDTOs
	var hi int = orderwList[0].UnitProfit
	itemStatis := new(dto.ItemStatisDTO)
	for i := 0; i < len(orderwList); i++ {
		orderw := orderwList[i]
		if orderw.UnitProfit >= hi-100 {
			itemStatis.UnitProfitRange = fmt.Sprintf("%d ~ %d", hi, orderw.UnitProfit)
			itemStatis.Orderwrappers = append(itemStatis.Orderwrappers, orderw)
			if i != len(orderwList)-1 {
				continue
			}
		}
		itemStatis.GenerateUnitProfit(unitLpCost)
		itemStatisDTOS = append(itemStatisDTOS, itemStatis)
		itemStatis = new(dto.ItemStatisDTO)
		hi = orderwList[i].UnitProfit
	}

	return &itemStatisDTOS, nil
}
