package model

import "time"

type ItemHistory struct {
	ItemId     int
	Average    float64   `json:"average"`
	Highest    float64   `json:"highest"`
	Lowest     float64   `json:"lowest"`
	OrderCount int64     `json:"order_count"`
	Volume     int64     `json:"volume"`
	Date       time.Time `json:"date"`
}

type ItemHistorys []*ItemHistory

func (ih ItemHistorys) Len() int { return len(ih) }

func (ih ItemHistorys) Less(i, j int) bool { return ih[i].ItemId < ih[j].ItemId }

func (ih ItemHistorys) Swap(i, j int) { ih[i], ih[j] = ih[j], ih[i] }
