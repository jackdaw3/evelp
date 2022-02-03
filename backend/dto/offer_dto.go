package dto

type OfferDTO struct {
	ItemId       int
	Name         string
	Quantity     int
	IskCost      float64
	LpCost       int
	Matertials   Matertials
	MaterialCost float64
	Price        float64
	Income       float64
	Profit       float64
	Volume       float64
	UnitProfit   int
	SalaIndex    int
	IsBluePrint  bool
}

type OfferDTOs []OfferDTO

func (o OfferDTOs) Len() int { return len(o) }

func (o OfferDTOs) Less(i, j int) bool { return o[i].UnitProfit > o[j].UnitProfit }

func (o OfferDTOs) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
