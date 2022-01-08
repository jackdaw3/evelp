package model

import (
	"database/sql/driver"
	"encoding/json"
	"evelp/configs/global"
	"fmt"
	"sort"
	"strconv"

	"gorm.io/gorm/clause"
)

type Offer struct {
	OfferId        int                `gorm:"type:int;not null;primary_key;autoIncrement:false" json:"offer_id"`
	ItemId         int                `gorm:"type:int;not null" json:"type_id"`
	Quantity       int                `gorm:"type:int;not null" json:"quantity"`
	IskCost        int                `gorm:"type:int;not null" json:"isk_cost"`
	LpCost         int                `gorm:"type:int;not null" json:"lp_cost"`
	AkCost         int                `gorm:"type:int;not null" json:"ak_cost"`
	RequireItems   RequireItems       `gorm:"type:text" json:"required_items"`
	CorporationIDs CorporationIDArray `gorm:"type:text"`
	IsBluePrint    bool               `gorm:"type:bool;default:false"`
}

type Offers []*Offer

type CorporationIDArray []int

func (offers Offers) Len() int { return len(offers) }

func (offers Offers) Less(i, j int) bool { return offers[i].OfferId < offers[j].OfferId }

func (offers Offers) Swap(i, j int) { offers[i], offers[j] = offers[j], offers[i] }

func (corporationIDs CorporationIDArray) Len() int { return len(corporationIDs) }

func (corporationIDs CorporationIDArray) Less(i, j int) bool {
	return corporationIDs[i] < corporationIDs[j]
}

func (corporationIDs CorporationIDArray) Swap(i, j int) {
	corporationIDs[i], corporationIDs[j] = corporationIDs[j], corporationIDs[i]
}

func (corporationIDs *CorporationIDArray) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}

	return json.Unmarshal(str, &corporationIDs)
}

func (corporationIDs CorporationIDArray) Value() (driver.Value, error) {
	str, err := json.Marshal(corporationIDs)
	if err != nil {
		return nil, nil
	}

	return string(str), nil
}

func GetOffer(offerId int) (*Offer, error) {
	var offer Offer
	result := global.DB.First(&offer, offerId)
	return &offer, result.Error
}

func GetCorporationOffer(corporationId int) (*Offers, error) {
	var offers Offers
	criteria := fmt.Sprintf("%%%s%%", strconv.Itoa(corporationId))
	result := global.DB.Where("corporation_ids LIKE ?", criteria).Find(&offers)
	return &offers, result.Error
}

func SaveOffer(offer *Offer) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&offer).Error; err != nil {
		return err
	}
	return nil
}

func SaveOffers(offers *Offers) error {
	for _, offer := range *offers {
		sort.Sort(offer.CorporationIDs)
		if err := SaveOffer(offer); err != nil {
			return err
		}
	}
	return nil
}
