package service

import (
	"evelp/dto"
	"evelp/model"
	"fmt"
	"math"
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

	orderService := NewOrderService(offer.ItemId, is.regionId, offer.IsBluePrint, is.scope)
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

	var (
		itemStatisDTOS dto.ItemStatisDTOs
		deviation      int
	)
	itemStatis := new(dto.ItemStatisDTO)

	for i := 0; i < len(orderwList); {
		tmp := orderwList[i].UnitProfit
		for i < len(orderwList) {
			deviation = int(math.Abs(float64(tmp - orderwList[i].UnitProfit)))
			if deviation > 100 {
				break
			}
			itemStatis.Orderwrappers = append(itemStatis.Orderwrappers, orderwList[i])
			i++
		}

		itemStatis.UnitProfitRange = fmt.Sprintf("%d ~ %d", tmp, orderwList[i-1].UnitProfit)
		itemStatis.GenerateUnitProfit(unitLpCost)
		itemStatisDTOS = append(itemStatisDTOS, itemStatis)

		itemStatis = new(dto.ItemStatisDTO)
		deviation = 0
	}

	return &itemStatisDTOS, nil
}
