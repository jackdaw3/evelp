package service

import (
	"evelp/dto"
	"evelp/model"
	"evelp/util/language"
	"fmt"
	"sort"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type OfferSerivce struct {
	regionId int
	scope    float64
	lang     string
}

func NewOfferSerivce(regionId int, scope float64, lang string) *OfferSerivce {
	return &OfferSerivce{regionId, scope, lang}
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
			offerDTO, err = o.convertBluePrint(offer)

		} else {
			offerDTO, err = o.convertOffer(offer)
		}

		if err != nil {
			return nil, err
		}
		offerDTOs = append(offerDTOs, *offerDTO)
	}

	sort.Sort(offerDTOs)

	return &offerDTOs, nil

}

func (o *OfferSerivce) convertOffer(offer *model.Offer) (*dto.OfferDTO, error) {
	var offerDTO dto.OfferDTO

	item, err := model.GetItem(offer.ItemId)
	if err != nil {
		return nil, err
	}

	offerDTO.ItemId = item.ItemId
	offerDTO.Name = language.Name(o.lang, item.Name)
	offerDTO.IsBluePrint = false
	offerDTO.Quantity = offer.Quantity
	offerDTO.IskCost = offer.IskCost
	offerDTO.LpCost = offer.LpCost

	materails, err := o.conertMaterials(offer.RequireItems)
	if err != nil {
		return nil, errors.WithMessage(err, "covert materails failed")
	}

	offerDTO.Matertials = materails
	offerDTO.MaterialCost = materails.Cost()

	oos := NewOrderService(offerDTO.ItemId, o.regionId, o.scope)
	price, err := oos.HighestBuyPrice()
	if err != nil {
		log.Errorf(err.Error())
	}
	offerDTO.Price = price
	offerDTO.Income = offerDTO.Price * float64(offer.Quantity)
	offerDTO.Profit = offerDTO.Income - (offerDTO.MaterialCost + offerDTO.IskCost)

	if offerDTO.LpCost > 0 {
		offerDTO.UnitProfit = int(offerDTO.Profit / float64(offerDTO.LpCost))
	}

	//TODO SET SALE INDEX
	return &offerDTO, nil
}

func (o *OfferSerivce) convertBluePrint(offer *model.Offer) (*dto.OfferDTO, error) {
	var offerDTO dto.OfferDTO

	bluePrint := model.GetBluePrint(offer.ItemId)
	if len(bluePrint.Products) == 0 {
		return nil, fmt.Errorf("offer %d's bluePrint %d have no product", offer.OfferId, bluePrint.BlueprintId)
	}

	bluePrintItem, err := model.GetItem(bluePrint.BlueprintId)
	if err != nil {
		return nil, err
	}

	product, err := model.GetItem(bluePrint.Products[0].ItemId)
	if err != nil {
		return nil, err
	}

	offerDTO.ItemId = product.ItemId
	offerDTO.Name = language.Name(o.lang, bluePrintItem.Name)
	offerDTO.IsBluePrint = true
	offerDTO.Quantity = offer.Quantity
	offerDTO.IskCost = offer.IskCost
	offerDTO.LpCost = offer.LpCost

	materails, err := o.conertMaterials(offer.RequireItems)
	if err != nil {
		return nil, errors.WithMessage(err, "covert materails failed")
	}
	manufactMaterials, err := o.conertManufactMaterials(bluePrint.Materials)
	if err != nil {
		return nil, errors.WithMessage(err, "covert manufact materials failed")
	}
	materails = append(materails, manufactMaterials...)

	offerDTO.Matertials = materails
	offerDTO.MaterialCost = materails.Cost()

	oos := NewOrderService(offerDTO.ItemId, o.regionId, o.scope)
	price, err := oos.HighestBuyPrice()
	if err != nil {
		log.Errorf(err.Error())
	}
	offerDTO.Price = price
	offerDTO.Income = offerDTO.Price * float64(offer.Quantity)
	offerDTO.Profit = offerDTO.Income - (offerDTO.MaterialCost + offerDTO.IskCost)

	if offerDTO.LpCost > 0 {
		offerDTO.UnitProfit = int(offerDTO.Profit / float64(offerDTO.LpCost))
	}

	//TODO SET SALE INDEX
	return &offerDTO, nil
}

func (o *OfferSerivce) conertMaterials(rs model.RequireItems) (dto.Matertials, error) {
	var materails dto.Matertials

	for _, r := range rs {
		var materail dto.Material
		mi, err := model.GetItem(r.ItemId)
		if err != nil {
			return nil, err
		}

		materail.ItemId = mi.ItemId
		materail.Name = language.Name(o.lang, mi.Name)

		materail.Quantity = r.Quantity
		materail.IsBluePrint = false

		mos := NewOrderService(mi.ItemId, o.regionId, o.scope)
		price, err := mos.LowestSellPrice()
		if err != nil {
			return nil, err
		}
		materail.Price = price
		materail.Cost = materail.Price * float64(materail.Quantity)
		materails = append(materails, materail)
	}

	return materails, nil
}

func (o *OfferSerivce) conertManufactMaterials(ms model.ManufactMaterials) (dto.Matertials, error) {
	var materails dto.Matertials

	for _, m := range ms {
		var materail dto.Material
		mi, err := model.GetItem(m.ItemId)
		if err != nil {
			return nil, err
		}

		materail.ItemId = mi.ItemId
		materail.Name = language.Name(o.lang, mi.Name)
		materail.IsBluePrint = true
		materail.Quantity = m.Quantity

		mos := NewOrderService(mi.ItemId, o.regionId, o.scope)
		price, err := mos.LowestSellPrice()
		if err != nil {
			return nil, err
		}
		materail.Price = price
		materail.Cost = materail.Price * float64(materail.Quantity)

		materails = append(materails, materail)
	}

	return materails, nil
}
