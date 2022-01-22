package dto

import "evelp/model"

type OffersDTO struct {
	Item                model.Item
	Quantity            int64
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
