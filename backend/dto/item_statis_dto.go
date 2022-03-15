package dto

type ItemStatisDTO struct {
	UnitProfitRange string
	AveUnitProfit   int
	Quantity        int64
	Income          int64
	Cost            int64
	Profit          int64
	Orderwrappers   OrderDTOWrappers
}

type ItemStatisDTOs []*ItemStatisDTO

func (is *ItemStatisDTO) GenerateUnitProfit(unitLpCost int) {
	var (
		quantity int64
		income   float64
		cost     float64
		profit   float64
	)

	for _, orderw := range is.Orderwrappers {
		quantity += orderw.OrderDTO.VolumeRemain
		income += orderw.Income
		cost += orderw.Cost
		profit += orderw.Profit
	}

	is.Quantity = quantity
	is.Income = int64(income)
	is.Cost = int64(cost)
	is.Profit = int64(profit)
	is.AveUnitProfit = int(is.Profit / (int64(unitLpCost) * is.Quantity))
}

func (is ItemStatisDTOs) Len() int { return len(is) }

func (is ItemStatisDTOs) Less(i, j int) bool { return is[i].AveUnitProfit < is[j].AveUnitProfit }

func (is ItemStatisDTOs) Swap(i, j int) { is[i], is[j] = is[j], is[i] }

type OrderDTOWrapper struct {
	OrderDTO   OrderDTO
	Income     float64
	Cost       float64
	Profit     float64
	UnitProfit int
}

type OrderDTOWrappers []*OrderDTOWrapper

func (o OrderDTOWrappers) Len() int { return len(o) }

func (o OrderDTOWrappers) Less(i, j int) bool { return o[i].UnitProfit < o[j].UnitProfit }

func (o OrderDTOWrappers) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
