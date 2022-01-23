package service

import (
	"evelp/dto"
	"evelp/model"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
)

type OfferSerivce struct {
	regionId int
	scope    float64
}

func NewOfferSerivce(regionId int, scope float64) *OfferSerivce {
	return &OfferSerivce{regionId, scope}
}

func (o *OfferSerivce) Offers(corporationId int) (*dto.OfferDTOs, error) {
	offers, err := model.GetOffersByCorporation(corporationId)
	if err != nil {
		return nil, err
	}

	var offerDTOs dto.OfferDTOs
	for _, offer := range *offers {
		var offerDTO *dto.OfferDTO
		var err error
		if offer.IsBluePrint {
			offerDTO, err = convertBluePrint(offer, o.regionId, o.scope)

		} else {
			offerDTO, err = convertOffer(offer, o.regionId, o.scope)
		}

		if err != nil {
			return nil, err
		}
		offerDTOs = append(offerDTOs, *offerDTO)
	}

	sort.Sort(offerDTOs)

	return &offerDTOs, nil

}

func convertOffer(offer *model.Offer, regionId int, scope float64) (*dto.OfferDTO, error) {
	var offerDTO dto.OfferDTO

	item, err := model.GetItem(offer.ItemId)
	if err != nil {
		return nil, err
	}
	offerDTO.Item = *item
	offerDTO.Quantity = offer.Quantity
	offerDTO.IskCost = offer.IskCost
	offerDTO.LpCost = offer.LpCost

	var materails dto.Matertials
	requireItems := offer.RequireItems
	for _, r := range requireItems {
		var materail dto.Material
		mi, err := model.GetItem(r.ItemId)
		if err != nil {
			return nil, err
		}
		materail.Item = *mi
		materail.Quantity = r.Quantity
		materail.IsBluePrint = false

		mos := NewOrderService(mi.ItemId, regionId, scope)
		price, err := mos.LowestSellPrice()
		if err != nil {
			return nil, err
		}
		materail.Price = price
		materails = append(materails, materail)
	}

	offerDTO.Matertials = materails
	offerDTO.MaterialCost = materails.Cost()

	oos := NewOrderService(offerDTO.Item.ItemId, regionId, scope)
	price, err := oos.HighestBuyPrice()
	if err != nil {
		log.Errorf(err.Error())
	}
	offerDTO.Income = price * float64(offerDTO.Quantity)
	offerDTO.Profit = offerDTO.Income - (offerDTO.MaterialCost + offerDTO.IskCost)

	if offerDTO.LpCost > 0 {
		offerDTO.LoyaltyPointsPerIsk = int(offerDTO.Profit / float64(offerDTO.LpCost))
	}
	//TODO SET SALE INDEX

	return &offerDTO, nil
}

func convertBluePrint(offer *model.Offer, regionId int, scope float64) (*dto.OfferDTO, error) {
	var offerDTO dto.OfferDTO

	bluePrint := model.GetBluePrint(offer.ItemId)
	if len(bluePrint.Products) == 0 {
		return nil, fmt.Errorf("offer %d's bluePrint %d have no product", offer.OfferId, bluePrint.BlueprintId)
	}

	item, err := model.GetItem(bluePrint.Products[0].ItemId)
	if err != nil {
		return nil, err
	}
	offerDTO.Item = *item
	offerDTO.Quantity = offer.Quantity
	offerDTO.IskCost = offer.IskCost
	offerDTO.LpCost = offer.LpCost

	var materails dto.Matertials
	manufactMaterials := bluePrint.Materials
	for _, m := range manufactMaterials {
		var materail dto.Material
		mi, err := model.GetItem(m.ItemId)
		if err != nil {
			return nil, err
		}
		materail.Item = *mi
		materail.IsBluePrint = true
		materail.Quantity = int64(m.Quantity)

		mos := NewOrderService(mi.ItemId, regionId, scope)
		price, err := mos.LowestSellPrice()
		if err != nil {
			return nil, err
		}
		materail.Price = price

		materails = append(materails, materail)
	}

	offerDTO.Matertials = materails
	offerDTO.MaterialCost = materails.Cost()

	oos := NewOrderService(offerDTO.Item.ItemId, regionId, scope)
	price, err := oos.HighestBuyPrice()
	if err != nil {
		log.Errorf(err.Error())
	}
	offerDTO.Income = price
	offerDTO.Profit = offerDTO.Income - (offerDTO.MaterialCost + offerDTO.IskCost)

	if offerDTO.LpCost > 0 {
		offerDTO.LoyaltyPointsPerIsk = int(offerDTO.Profit / float64(offerDTO.LpCost))
	}
	//TODO SET SALE INDEX

	return &offerDTO, nil
}
