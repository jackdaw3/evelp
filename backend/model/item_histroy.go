package model

type ItemHistory struct {
	ItemId     int
	Average    float64
	Highest    float64
	Lowest     float64
	OrderCount int64
	Volume     int64
}

type ItemHistorys []*ItemHistory

func (ih ItemHistorys) Len() int { return len(ih) }

func (ih ItemHistorys) Less(i, j int) bool { return ih[i].ItemId < ih[j].ItemId }

func (ih ItemHistorys) Swap(i, j int) { ih[i], ih[j] = ih[j], ih[i] }
