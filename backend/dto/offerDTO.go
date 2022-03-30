package dto

import "math"

type OfferDTO struct {
	OfferId      int
	ItemId       int
	ItemName     string
	Quantity     int
	IskCost      float64
	LpCost       int
	Matertials   MatertialDTOs
	MaterialCost float64
	Price        float64
	Income       float64
	Profit       float64
	Volume       int64
	UnitProfit   int
	SaleIndex    int
	IsBluePrint  bool
	Error        bool
	ErrorMessage string
}

const (
	k1 = 100000000
	k2 = 100
	k3 = 45
)

type OfferDTOs []OfferDTO

func (o OfferDTOs) Len() int { return len(o) }

func (o OfferDTOs) Less(i, j int) bool { return o[i].UnitProfit > o[j].UnitProfit }

func (o OfferDTOs) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

func (o *OfferDTO) GenerateSaleIndex() {
	var (
		saleIndex       int
		quantityScore   float64
		iskStreamScore  float64
		unitProfitScore float64
	)

	quantityScore = (2.25/math.Pi)*math.Atan(float64(o.Volume)/20) + 1
	iskStream := o.Price * float64(o.Volume)
	if iskStream <= k1 {
		iskStreamScore = iskStream / k1
	} else {
		iskStreamScore = (math.Log(iskStream/k1))/(math.Log(2)) + 1
	}
	iskStreamScore = iskStreamScore / 2
	unitProfitScore = float64(o.UnitProfit) / k2
	if unitProfitScore > k3 {
		unitProfitScore = k3
	}

	saleIndex = int(quantityScore * iskStreamScore * unitProfitScore)
	o.SaleIndex = saleIndex
}
