package dto

import "evelp/model"

type OfferDTO struct {
	Item                model.Item
	Quantity            int
	IskCost             float64
	LpCost              int
	Matertials          Matertials
	MaterialCost        float64
	Income              float64
	Profit              float64
	VolumePerDay        float64
	LoyaltyPointsPerIsk int
	SalaIndex           int
}

type OfferDTOs []OfferDTO

func (o OfferDTOs) Len() int { return len(o) }

func (o OfferDTOs) Less(i, j int) bool { return o[i].LoyaltyPointsPerIsk > o[j].LoyaltyPointsPerIsk }

func (o OfferDTOs) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
