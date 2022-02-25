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

	var orderService *OrderService
	if offer.IsBluePrint {
		bluePrint := model.GetBluePrint(offer.ItemId)
		product, err := model.GetItem(bluePrint.Products[0].ItemId)
		if err != nil {
			return nil, err
		}
		orderService = NewOrderService(product.ItemId, is.regionId, is.scope)
	} else {
		orderService = NewOrderService(offer.ItemId, is.regionId, is.scope)
	}
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
		var orderw dto.OrderDTOWrapper
		orderw.OrderDTO = order
		orderw.Income = order.Price * float64(order.VolumeRemain)
		orderw.Cost = unitCost * float64(order.VolumeRemain)
		orderw.Profit = orderw.Income - orderw.Cost
		orderw.UnitProfit = int(orderw.Profit / float64(offer.LpCost))
		orderwList = append(orderwList, &orderw)
	}
	sort.Sort(orderwList)

	var itemStatisDTOS dto.ItemStatisDTOs
	var hi int
	var itemStatis dto.ItemStatisDTO
	for i := 0; i < len(orderwList); {
		orderw := orderwList[i]
		if hi == 0 {
			hi = orderw.UnitProfit
			itemStatis.Orderwrappers = append(itemStatis.Orderwrappers, orderw)
		} else {
			if orderw.UnitProfit >= hi-100 {
				itemStatis.UnitProfitRange = fmt.Sprintf("%d ~ %d", hi, orderw.UnitProfit)
				itemStatis.Orderwrappers = append(itemStatis.Orderwrappers, orderw)
				i++
			} else {
				itemStatis.GenerateUnitProfit(unitLpCost)
				itemStatisDTOS = append(itemStatisDTOS, &itemStatis)
				itemStatis = dto.ItemStatisDTO{}
				hi = 0
			}
		}
	}

	return &itemStatisDTOS, nil
}
