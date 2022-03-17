package service

import (
	"context"
	"evelp/dto"
	"evelp/log"
	"evelp/model"
	"fmt"
	"sort"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"
)

const (
	weight = 1
	limit  = 3
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

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	sem := semaphore.NewWeighted(limit)

	var offerDTOs dto.OfferDTOs
	for _, offer := range *offers {
		sem.Acquire(context.Background(), weight)
		wg.Add(1)

		go func(offer *model.Offer) {
			defer sem.Release(weight)
			defer wg.Done()

			var (
				offerDTO *dto.OfferDTO
				err      error
			)

			if offer.IsBluePrint {
				offerDTO, err = o.convertBluePrint(offer)
			} else {
				offerDTO, err = o.convertOffer(offer)
			}

			if err != nil {
				log.Errorf(err, "get offer %d failed", offer.OfferId)
				return
			}

			defer mu.Unlock()
			mu.Lock()
			offerDTOs = append(offerDTOs, *offerDTO)
		}(offer)
	}

	wg.Wait()
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

	bluePrint, err := model.GetBluePrint(offer.ItemId)
	if err != nil {
		return nil, errors.WithMessagef(err, "get blue print %d failed", offer.ItemId)
	}
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
	var materials dto.MatertialDTOs

	for _, r := range rs {
		var material dto.MaterialDTO
		mi, err := model.GetItem(r.ItemId)
		if err != nil {
			log.Errorf(err, "get item %v failed", r.ItemId)
			continue
		}

		material.ItemId = mi.ItemId
		material.Name = mi.Name.Val(o.lang)

		material.Quantity = r.Quantity
		material.IsBluePrint = false

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
				material.Name,
				errors.Cause(err).Error(),
			)
			if len(offerDTO.ErrorMessage) > 0 {
				offerDTO.ErrorMessage += ".\n" + errorMessage
			} else {
				offerDTO.ErrorMessage = errorMessage
			}
			material.Error = true
			material.ErrorMessage = errorMessage
			log.Warnf("get %s price of item %v in region %v failed: %v", o.materialPrice, mos.itemId, mos.regionId, err)
		}
		material.Price = price
		material.Cost = material.Price * float64(material.Quantity)
		materials = append(materials, material)
	}

	return materials
}

func (o *OfferSerivce) conertManufactMaterials(ms model.ManufactMaterials, offerDTO *dto.OfferDTO) dto.MatertialDTOs {
	var materials dto.MatertialDTOs

	for _, m := range ms {
		var material dto.MaterialDTO
		mi, err := model.GetItem(m.ItemId)
		if err != nil {
			log.Errorf(err, "get item %v failed", m.ItemId)
			continue
		}

		material.ItemId = mi.ItemId
		material.Name = mi.Name.Val(o.lang)
		material.IsBluePrint = true
		material.Quantity = m.Quantity

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
				material.Name,
				errors.Cause(err).Error(),
			)
			if len(offerDTO.ErrorMessage) > 0 {
				offerDTO.ErrorMessage += ".\n" + errorMessage
			} else {
				offerDTO.ErrorMessage = errorMessage
			}
			material.Error = true
			material.ErrorMessage = errorMessage
			log.Warnf("get %s price of item %v in region %v failed: %v", o.materialPrice, mos.itemId, mos.regionId, err)
		}
		material.Price = price
		material.Cost = material.Price * float64(material.Quantity)

		materials = append(materials, material)
	}

	return materials
}
