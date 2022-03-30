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
	tax           float64
	lang          string
}

func NewItemStatisService(offerId, regionId int, scope float64, materialPrice string, tax float64, lang string) *ItemStatisService {
	return &ItemStatisService{offerId, regionId, scope, materialPrice, tax, lang}
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

	offerService := NewOfferSerivce(is.regionId, is.scope, 0, "", is.materialPrice, is.tax, is.lang)
	offerDTO, err := offerService.Offer(is.offerId)
	if err != nil {
		return nil, err
	}

	unitCost := (offerDTO.MaterialCost + offerDTO.IskCost) / float64(offerDTO.Quantity)
	unitLpCost := float64(offerDTO.LpCost) / float64(offerDTO.Quantity)

	var orderwList dto.OrderDTOWrappers
	for _, order := range *orders {
		orderw := new(dto.OrderDTOWrapper)
		orderw.OrderDTO = order

		if is.tax > 0 {
			orderw.Income = order.Price * ((100 - is.tax) / 100) * float64(order.VolumeRemain)
		} else {
			orderw.Income = order.Price * float64(order.VolumeRemain)
		}
		orderw.Cost = unitCost * float64(order.VolumeRemain)
		orderw.Profit = orderw.Income - orderw.Cost
		orderw.UnitProfit = int(orderw.Profit / (unitLpCost * float64(order.VolumeRemain)))
		orderwList = append(orderwList, orderw)
	}
	if isBuyOrder {
		sort.Sort(sort.Reverse(orderwList))
	} else {
		sort.Sort(orderwList)
	}

	var (
		itemStatisDTOs dto.ItemStatisDTOs
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
		itemStatisDTOs = append(itemStatisDTOs, itemStatis)

		itemStatis = new(dto.ItemStatisDTO)
		deviation = 0
	}

	return &itemStatisDTOs, nil
}
