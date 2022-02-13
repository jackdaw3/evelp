package dto

type ItemStatisDTO struct {
	UnitProfitRange  string
	AveUnitProfit    int
	Quantity         int64
	Income           float64
	Cost             float64
	Profit           float64
	OrderDTOWrappers OrderDTOWrappers
}

type ItemStatisDTOs []*ItemStatisDTO

func (is ItemStatisDTOs) Len() int { return len(is) }

func (is ItemStatisDTOs) Less(i, j int) bool { return is[i].AveUnitProfit > is[j].AveUnitProfit }

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

func (o OrderDTOWrappers) Less(i, j int) bool { return o[i].Profit > o[j].Profit }

func (o OrderDTOWrappers) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
