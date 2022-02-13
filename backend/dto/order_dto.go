package dto

import "time"

type OrderDTO struct {
	OrderId      int
	ItemId       int
	ItemName     string
	SystemName   string
	Issued       time.Time
	Duration     int
	VolumeRemain int64
	VolumeTotal  int64
	IsBuyOrder   bool

	Price       float64
	LastUpdated time.Time
}

type OrderDTOs []OrderDTO

func (o OrderDTOs) Len() int { return len(o) }

func (o OrderDTOs) Less(i, j int) bool { return o[i].Price < o[j].Price }

func (o OrderDTOs) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
