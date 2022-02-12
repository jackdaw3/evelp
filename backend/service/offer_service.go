package service

import (
	"evelp/dto"
	"evelp/log"
	"evelp/model"
	"sort"

	"github.com/pkg/errors"
)

type OfferSerivce struct {
	regionId int
	scope    float64
	days     int
	lang     string
}

func NewOfferSerivce(regionId int, scope float64, days int, lang string) *OfferSerivce {
	return &OfferSerivce{regionId, scope, days, lang}
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

	offerDTO.ItemId = item.ItemId
	offerDTO.Name = item.Name.Val(o.lang)
	offerDTO.IsBluePrint = false
	offerDTO.Quantity = offer.Quantity
	offerDTO.IskCost = offer.IskCost
	offerDTO.LpCost = offer.LpCost

	materails := o.conertMaterials(offer.RequireItems)
	offerDTO.Matertials = materails
	offerDTO.MaterialCost = materails.Cost()

	oos := NewOrderService(offerDTO.ItemId, o.regionId, o.scope)
	price, err := oos.HighestBuyPrice()
	if err != nil {
		log.Warnf("get highest buy price of item %v in region %v failed: %v", oos.itemId, oos.regionId, err)
	}
	offerDTO.Price = price
	offerDTO.Income = offerDTO.Price * float64(offer.Quantity)
	offerDTO.Profit = offerDTO.Income - (offerDTO.MaterialCost + offerDTO.IskCost)

	if offerDTO.LpCost > 0 {
		offerDTO.UnitProfit = int(offerDTO.Profit / float64(offerDTO.LpCost))
	}

	volume, err := o.AverageVolume(offerDTO.ItemId)
	if err != nil {
		log.Warnf("get volume of item %v region %v failed: %v", oos.itemId, oos.regionId, err)
	}
	offerDTO.Volume = volume

	//TODO SET SALE INDEX
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

	product, err := model.GetItem(bluePrint.Products[0].ItemId)
	if err != nil {
		return nil, err
	}

	offerDTO.ItemId = product.ItemId
	offerDTO.IsBluePrint = true
	offerDTO.Quantity = offer.Quantity
	offerDTO.IskCost = offer.IskCost
	offerDTO.LpCost = offer.LpCost

	materails := o.conertMaterials(offer.RequireItems)
	manufactMaterials := o.conertManufactMaterials(bluePrint.Materials)
	materails = append(materails, manufactMaterials...)

	offerDTO.Matertials = materails
	offerDTO.MaterialCost = materails.Cost()

	oos := NewOrderService(offerDTO.ItemId, o.regionId, o.scope)
	price, err := oos.HighestBuyPrice()
	if err != nil {
		log.Warnf("get highest buy price of item %v in region %v failed: %v", oos.itemId, oos.regionId, err)
	}
	offerDTO.Price = price
	offerDTO.Income = offerDTO.Price * float64(offer.Quantity)
	offerDTO.Profit = offerDTO.Income - (offerDTO.MaterialCost + offerDTO.IskCost)

	if offerDTO.LpCost > 0 {
		offerDTO.UnitProfit = int(offerDTO.Profit / float64(offerDTO.LpCost))
	}

	volume, err := o.AverageVolume(offerDTO.ItemId)
	if err != nil {
		log.Warnf("get volume of item %v region %v failed: %v", oos.itemId, oos.regionId, err)
	}
	offerDTO.Volume = volume

	//TODO SET SALE INDEX

	offerDTO.ItemId = bluePrintItem.ItemId
	offerDTO.Name = bluePrintItem.Name.Val(o.lang)

	return &offerDTO, nil
}

func (o *OfferSerivce) conertMaterials(rs model.RequireItems) dto.MatertialDTOs {
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

		mos := NewOrderService(mi.ItemId, o.regionId, o.scope)
		price, err := mos.LowestSellPrice()
		if err != nil {
			log.Warnf("get lowest sell price of item %v in region %v failed: %v", mos.itemId, mos.regionId, err)
		}
		materail.Price = price
		materail.Cost = materail.Price * float64(materail.Quantity)
		materails = append(materails, materail)
	}

	return materails
}

func (o *OfferSerivce) conertManufactMaterials(ms model.ManufactMaterials) dto.MatertialDTOs {
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

		mos := NewOrderService(mi.ItemId, o.regionId, o.scope)
		price, err := mos.LowestSellPrice()
		if err != nil {
			log.Warnf("get lowest sell price of item %v in region %v failed: %v", mos.itemId, mos.regionId, err)
		}
		materail.Price = price
		materail.Cost = materail.Price * float64(materail.Quantity)

		materails = append(materails, materail)
	}

	return materails
}

func (o *OfferSerivce) AverageVolume(itemId int) (int64, error) {
	ihs := NewItemHistoryService(itemId, o.regionId)
	itemHistorys, err := ihs.History()
	if err != nil {
		return 0, err
	}

	volume := itemHistorys.AverageVolume(o.days)
	return volume, nil
}
