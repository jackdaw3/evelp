package dto

import "evelp/model"

type Material struct {
	Item        model.Item
	Quantity    int64
	Price       float64
	IsBluePrint bool
}

type Matertials []Material

func (ms *Matertials) Cost() float64 {
	var cost float64
	for _, m := range *ms {
		cost += m.Price * float64(m.Quantity)
	}
	return cost
}
