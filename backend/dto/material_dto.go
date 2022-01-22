package dto

import "evelp/model"

type Material struct {
	Item        model.Item
	Quantity    int64
	Price       float64
	IsBluePrint bool
}

type Matertials []Material
