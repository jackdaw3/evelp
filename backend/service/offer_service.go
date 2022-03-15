package service

import (
	"evelp/dto"
	"evelp/log"
	"evelp/model"
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

type OfferSerivce struct {
	regionId      int
	scope         float64
	days          int
	productPrice  string
	materialPrice string
	tax           float64
	lang          string
}

func NewOfferSerivce(regionId int, scope float64, days int, productPrice string, materialPrice string, tax float64, lang string) *OfferSerivce {
	return &OfferSerivce{regionId, scope, days, productPrice, materialPrice, tax, lang}
}

func (o *OfferSerivce) Offer(offerId int) (*dto.OfferDTO, error) {
	offer, err := model.GetOffer(offerId)
	if err != nil {
		return nil, err
	}

	var offerDTO *dto.OfferDTO
	if offer.IsBluePrint {
		offerDTO, err = o.convertBluePrint(offer)
	} else {
		offerDTO, err = o.convertOffer(offer)
	}

	if err != nil {
		return nil, errors.WithMessagef(err, "get offer %d failed", offer.OfferId)
	}

	return offerDTO, nil
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
			log.Errorf(err, "get offer %d failed", offer.OfferId)
			continue
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

	offerDTO.OfferId = offer.OfferId
	offerDTO.ItemId = item.ItemId
	offerDTO.Name = item.Name.Val(o.lang)
	offerDTO.IsBluePrint = false
	offerDTO.Quantity = offer.Quantity
	offerDTO.IskCost = offer.IskCost
	offerDTO.LpCost = offer.LpCost

	materails := o.conertMaterials(offer.RequireItems, &offerDTO)
	offerDTO.Matertials = materails
	offerDTO.MaterialCost = materails.Cost()

	oos := NewOrderService(offerDTO.ItemId, o.regionId, false, o.scope)
	var price float64
	if o.productPrice == "buy" {
		price, err = oos.HighestBuyPrice()
	} else if o.productPrice == "sell" {
		price, err = oos.LowestSellPrice()
	}
	if err != nil {
		offerDTO.Error = true
		errorMessage := fmt.Sprintf("get %s price of product %s in The Forge failed: %s",
			o.productPrice, offerDTO.Name,
			errors.Cause(err).Error(),
		)
		if len(offerDTO.ErrorMessage) > 0 {
			offerDTO.ErrorMessage += ".\n" + errorMessage
		} else {
			offerDTO.ErrorMessage = errorMessage
		}

		log.Warnf("get %s price of item %v in region %v failed: %v", o.productPrice, oos.itemId, oos.regionId, err)
	}
	offerDTO.Price = price
	offerDTO.Income = offerDTO.Price * ((100 - o.tax) / 100) * float64(offer.Quantity)
	offerDTO.Profit = offerDTO.Income - (offerDTO.MaterialCost + offerDTO.IskCost)

	if offerDTO.LpCost > 0 {
		offerDTO.UnitProfit = int(offerDTO.Profit / float64(offerDTO.LpCost))
	}

	ihs := NewItemHistoryService(offerDTO.ItemId, o.regionId, offerDTO.IsBluePrint)
	volume, err := ihs.AverageVolume(o.days)
	if err != nil {
		log.Warnf("get volume of item %v region %v failed: %v", oos.itemId, oos.regionId, err)
	}
	offerDTO.Volume = volume
	offerDTO.GenerateSaleIndex()

	return &offerDTO, nil
}

func (o *OfferSerivce) convertBluePrint(offer *model.Offer) (*dto.OfferDTO, error) {
	var offerDTO dto.OfferDTO

	bluePrint := model.GetBluePrint(offer.ItemId)
	if len(bluePrint.Products) == 0 {
		return nil, errors.Errorf("offer %d's bluePrint %d have no product", offer.OfferId, bluePrint.BlueprintId)
	}

	bluePrintItem, err := model.GetItem(bluePrint.BlueprintId)
	if err != nil {
		return nil, err
	}

	offerDTO.OfferId = offer.OfferId
	offerDTO.ItemId = offer.ItemId
	offerDTO.IsBluePrint = true
	offerDTO.Quantity = offer.Quantity
	offerDTO.IskCost = offer.IskCost
	offerDTO.LpCost = offer.LpCost

	materails := o.conertMaterials(offer.RequireItems, &offerDTO)
	manufactMaterials := o.conertManufactMaterials(bluePrint.Materials, &offerDTO)
	materails = append(materails, manufactMaterials...)

	offerDTO.Matertials = materails
	offerDTO.MaterialCost = materails.Cost()

	oos := NewOrderService(offerDTO.ItemId, o.regionId, offerDTO.IsBluePrint, o.scope)
	var price float64
	if o.productPrice == "buy" {
		price, err = oos.HighestBuyPrice()
	} else if o.productPrice == "sell" {
		price, err = oos.LowestSellPrice()
	}
	if err != nil {
		offerDTO.Error = true
		errorMessage := fmt.Sprintf("get %s price of blueprint %s product in The Forge failed: %s",
			o.productPrice,
			bluePrintItem.Name.Val(o.lang),
			errors.Cause(err).Error(),
		)
		if len(offerDTO.ErrorMessage) > 0 {
			offerDTO.ErrorMessage += ".\n" + errorMessage
		} else {
			offerDTO.ErrorMessage = errorMessage
		}
		log.Warnf("get %s price of item %v in region %v failed: %v", o.productPrice, oos.itemId, oos.regionId, err)
	}
	offerDTO.Price = price
	offerDTO.Income = offerDTO.Price * ((100 - o.tax) / 100) * float64(offer.Quantity)
	offerDTO.Profit = offerDTO.Income - (offerDTO.MaterialCost + offerDTO.IskCost)

	if offerDTO.LpCost > 0 {
		offerDTO.UnitProfit = int(offerDTO.Profit / float64(offerDTO.LpCost))
	}

	ihs := NewItemHistoryService(offerDTO.ItemId, o.regionId, offerDTO.IsBluePrint)
	volume, err := ihs.AverageVolume(o.days)
	if err != nil {
		log.Warnf("get volume of item %v region %v failed: %v", oos.itemId, oos.regionId, err)
	}
	offerDTO.Volume = volume
	offerDTO.GenerateSaleIndex()
	offerDTO.Name = bluePrintItem.Name.Val(o.lang)

	return &offerDTO, nil
}

func (o *OfferSerivce) conertMaterials(rs model.RequireItems, offerDTO *dto.OfferDTO) dto.MatertialDTOs {
	var materails dto.MatertialDTOs

	for _, r := range rs {
		var materail dto.MaterialDTO
		mi, err := model.GetItem(r.ItemId)
		if err != nil {
			log.Errorf(err, "get item %v failed", r.ItemId)
			continue
		}

		materail.ItemId = mi.ItemId
		materail.Name = mi.Name.Val(o.lang)

		materail.Quantity = r.Quantity
		materail.IsBluePrint = false

		mos := NewOrderService(mi.ItemId, o.regionId, false, o.scope)
		var price float64
		if o.materialPrice == "sell" {
			price, err = mos.LowestSellPrice()
		} else if o.materialPrice == "buy" {
			price, err = mos.HighestBuyPrice()
		}
		if err != nil {
			offerDTO.Error = true
			errorMessage := fmt.Sprintf("get %s price of lp store material %s in The Forge failed: %s",
				o.materialPrice,
				materail.Name,
				errors.Cause(err).Error(),
			)
			if len(offerDTO.ErrorMessage) > 0 {
				offerDTO.ErrorMessage += ".\n" + errorMessage
			} else {
				offerDTO.ErrorMessage = errorMessage
			}
			materail.Error = true
			materail.ErrorMessage = errorMessage
			log.Warnf("get %s price of item %v in region %v failed: %v", o.materialPrice, mos.itemId, mos.regionId, err)
		}
		materail.Price = price
		materail.Cost = materail.Price * float64(materail.Quantity)
		materails = append(materails, materail)
	}

	return materails
}

func (o *OfferSerivce) conertManufactMaterials(ms model.ManufactMaterials, offerDTO *dto.OfferDTO) dto.MatertialDTOs {
	var materails dto.MatertialDTOs

	for _, m := range ms {
		var materail dto.MaterialDTO
		mi, err := model.GetItem(m.ItemId)
		if err != nil {
			log.Errorf(err, "get item %v failed", m.ItemId)
			continue
		}

		materail.ItemId = mi.ItemId
		materail.Name = mi.Name.Val(o.lang)
		materail.IsBluePrint = true
		materail.Quantity = m.Quantity

		mos := NewOrderService(mi.ItemId, o.regionId, false, o.scope)
		var price float64
		if o.materialPrice == "sell" {
			price, err = mos.LowestSellPrice()
		} else if o.materialPrice == "buy" {
			price, err = mos.HighestBuyPrice()
		}
		if err != nil {
			offerDTO.Error = true
			errorMessage := fmt.Sprintf("get %s price of manufact material %s in The Forge failed: %s",
				o.materialPrice,
				materail.Name,
				errors.Cause(err).Error(),
			)
			if len(offerDTO.ErrorMessage) > 0 {
				offerDTO.ErrorMessage += ".\n" + errorMessage
			} else {
				offerDTO.ErrorMessage = errorMessage
			}
			materail.Error = true
			materail.ErrorMessage = errorMessage
			log.Warnf("get %s price of item %v in region %v failed: %v", o.materialPrice, mos.itemId, mos.regionId, err)
		}
		materail.Price = price
		materail.Cost = materail.Price * float64(materail.Quantity)

		materails = append(materails, materail)
	}

	return materails
}
