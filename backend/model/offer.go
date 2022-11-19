package model

import (
	"database/sql/driver"
	"encoding/json"
	"evelp/config/global"
	"evelp/log"
	"evelp/util/cache"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

const (
	offer_key             = "offer"
	offers_key            = "offer:offers"
	offer_corporation_key = "offer:corporation"
)

type Offer struct {
	OfferId        int            `gorm:"type:int;not null;primary_key;autoIncrement:false" json:"offer_id"`
	ItemId         int            `gorm:"type:int;not null" json:"type_id"`
	Quantity       int            `gorm:"type:int;not null" json:"quantity"`
	IskCost        float64        `gorm:"type:double;not null" json:"isk_cost"`
	LpCost         int            `gorm:"type:int;not null" json:"lp_cost"`
	AkCost         int            `gorm:"type:int;not null" json:"ak_cost"`
	RequireItems   RequireItems   `gorm:"type:text" json:"required_items"`
	CorporationIds CorporationIds `gorm:"type:text" json:"corporation_ids"`
	IsBluePrint    bool           `gorm:"type:bool;default:false" json:"is_blue_print"`
}

type RequireItem struct {
	ItemId   int   `gorm:"type:int;not null" json:"type_id"`
	Quantity int64 `gorm:"type:int;not null" json:"quantity"`
}

type Offers []*Offer

type CorporationIds []int

type RequireItems []RequireItem

func (o Offers) Len() int { return len(o) }

func (o Offers) Less(i, j int) bool { return o[i].OfferId < o[j].OfferId }

func (o Offers) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

func (c CorporationIds) Len() int { return len(c) }

func (c CorporationIds) Less(i, j int) bool {
	return c[i] < c[j]
}

func (c CorporationIds) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c *CorporationIds) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return errors.Errorf("%v is not []byte", value)
	}

	return json.Unmarshal(str, &c)
}

func (c CorporationIds) Value() (driver.Value, error) {
	str, err := json.Marshal(c)
	if err != nil {
		return nil, nil
	}

	return string(str), nil
}

func (r *RequireItems) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return errors.Errorf("%v is not []byte", value)
	}

	return json.Unmarshal(str, &r)
}

func (r RequireItems) Value() (driver.Value, error) {
	str, err := json.Marshal(r)
	if err != nil {
		return nil, nil
	}

	return string(str), nil
}

func GetOffer(offerId int) (*Offer, error) {
	var offer Offer

	key := cache.Key(offer_key, strconv.Itoa(offerId))
	if err := cache.Get(key, &offer); err != nil {
		log.Debugf("failed to get offer %d from cache: %s", offerId, err.Error())
		result := global.DB.First(&offer, offerId)
		if err := cache.Set(key, offer, global.Conf.Redis.ExpireTime.Model*time.Minute); err != nil {
			return nil, err
		}
		return &offer, result.Error
	}

	return &offer, nil
}

func GetOffers() (*Offers, error) {
	var offers Offers

	if err := cache.Get(offers_key, &offers); err != nil {
		log.Debugf("failed to get all offers from cache: %s", err.Error())
		result := global.DB.Find(&offers)
		if err := cache.Set(offers_key, offers, global.Conf.Redis.ExpireTime.Model*time.Minute); err != nil {
			return nil, err
		}
		return &offers, result.Error
	}

	return &offers, nil
}

func GetOffersByCorporation(corporationId int) (*Offers, error) {
	var offers Offers

	key := cache.Key(offer_corporation_key, strconv.Itoa(corporationId))
	if err := cache.Get(key, &offers); err != nil {
		log.Debugf("failed to get corporation %d's offers from cache: %s", corporationId, err.Error())
		criteria := fmt.Sprintf("%%%s%%", strconv.Itoa(corporationId))
		result := global.DB.Where("corporation_ids LIKE ?", criteria).Find(&offers)
		if err := cache.Set(key, offers, global.Conf.Redis.ExpireTime.Model*time.Minute); err != nil {
			return nil, err
		}
		return &offers, result.Error
	}

	return &offers, nil
}

func SaveOffer(offer *Offer) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&offer).Error; err != nil {
		return errors.Wrapf(err, "failed to save offer %d to DB", offer.ItemId)
	}

	return nil
}

func SaveOffers(offers *Offers) error {
	for _, offer := range *offers {
		sort.Sort(offer.CorporationIds)
		if err := SaveOffer(offer); err != nil {
			return err
		}
	}
	return nil
}

func GetAllItems() (map[int]interface{}, error) {
	offers, err := GetOffers()
	if err != nil {
		return nil, nil
	}

	items := make(map[int]interface{})
	for _, offer := range *offers {
		if offer.IsBluePrint {
			bluePrint, err := GetBluePrint(offer.ItemId)
			if err != nil {
				log.Errorf(err, "failed to get blueprint %d", offer.ItemId)
			}
			if len(bluePrint.Products) > 0 {
				for _, product := range bluePrint.Products {
					items[product.ItemId] = struct{}{}
				}
			}
			if len(bluePrint.Materials) > 0 {
				for _, material := range bluePrint.Materials {
					items[material.ItemId] = struct{}{}
				}
			}

		} else {
			items[offer.ItemId] = struct{}{}
		}

		if len(offer.RequireItems) > 0 {
			for _, requireItem := range offer.RequireItems {
				items[requireItem.ItemId] = struct{}{}
			}
		}
	}

	return items, nil
}

func GetAllProducts() (map[int]interface{}, error) {
	offers, err := GetOffers()
	if err != nil {
		return nil, nil
	}

	items := make(map[int]interface{})
	for _, offer := range *offers {
		if offer.IsBluePrint {
			bluePrint, err := GetBluePrint(offer.ItemId)
			if err != nil {
				log.Errorf(err, "failed to get blueprint %d", offer.ItemId)
			}
			if len(bluePrint.Products) > 0 {
				for _, product := range bluePrint.Products {
					items[product.ItemId] = struct{}{}
				}
			}
		} else {
			items[offer.ItemId] = struct{}{}
		}
	}

	return items, nil
}
