package dto

type MaterialDTO struct {
	ItemId      int
	Name        string
	Quantity    int64
	Price       float64
	Cost        float64
	IsBluePrint bool
}

type MatertialDTOs []MaterialDTO

func (ms *MatertialDTOs) Cost() float64 {
	var cost float64
	for _, m := range *ms {
		cost += m.Cost
	}
	return cost
}
