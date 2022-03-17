package model

import (
	"database/sql/driver"
	"encoding/json"
	"evelp/config/global"
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
	offer_expiration      = -1 * time.Second
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

type Offers []*Offer

type CorporationIds []int

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

func GetOffer(offerId int) (*Offer, error) {
	var offer Offer

	key := cache.Key(offer_key, strconv.Itoa(offerId))
	exist := cache.Exist(key)

	if exist == nil {
		if err := cache.Get(key, &offer); err != nil {
			return nil, err
		}
		return &offer, nil
	} else {
		result := global.DB.First(&offer, offerId)
		if err := cache.Set(key, offer, offer_expiration); err != nil {
			return nil, err
		}
		return &offer, result.Error
	}
}

func GetOffers() (*Offers, error) {
	var offers Offers

	exist := cache.Exist(offers_key)

	if exist == nil {
		if err := cache.Get(offers_key, &offers); err != nil {
			return nil, err
		}
		return &offers, nil
	} else {
		result := global.DB.Find(&offers)
		if err := cache.Set(offers_key, offers, offer_expiration); err != nil {
			return nil, err
		}
		return &offers, result.Error
	}
}

func GetOffersByCorporation(corporationId int) (*Offers, error) {
	var offers Offers

	key := cache.Key(offer_corporation_key, strconv.Itoa(corporationId))
	exist := cache.Exist(key)

	if exist == nil {
		if err := cache.Get(key, &offers); err != nil {
			return nil, err
		}
		return &offers, nil
	} else {
		criteria := fmt.Sprintf("%%%s%%", strconv.Itoa(corporationId))
		result := global.DB.Where("corporation_ids LIKE ?", criteria).Find(&offers)
		if err := cache.Set(key, offers, offer_expiration); err != nil {
			return nil, err
		}
		return &offers, result.Error
	}
}

func SaveOffer(offer *Offer) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&offer).Error; err != nil {
		return errors.Wrap(err, "save offer to DB failed")
	}

	key := cache.Key(offer_key, strconv.Itoa(offer.OfferId))
	if err := cache.Set(key, *offer, offer_expiration); err != nil {
		return err
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
