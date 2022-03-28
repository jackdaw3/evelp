package model

type ItemHistory struct {
	ItemId     int
	Average    float64 `json:"average"`
	Highest    float64 `json:"highest"`
	Lowest     float64 `json:"lowest"`
	OrderCount int64   `json:"order_count"`
	Volume     int64   `json:"volume"`
	Date       string  `json:"date"`
}

type ItemHistorys []*ItemHistory

func (ih ItemHistorys) Len() int { return len(ih) }

func (ih ItemHistorys) Less(i, j int) bool { return ih[i].ItemId < ih[j].ItemId }

func (ih ItemHistorys) Swap(i, j int) { ih[i], ih[j] = ih[j], ih[i] }

func (i *ItemHistorys) AverageVolume(days int) int64 {
	size := len(*i)
	if size == 0 || days == 0 {
		return 0
	}
	var sum int64
	var count int64

	if size < days {
		for _, itemHistory := range *i {
			sum += itemHistory.Volume
			count++
		}
	} else {
		itemHistorys := (*i)[size-days:]
		for _, itemHistory := range itemHistorys {
			sum += itemHistory.Volume
			count++
		}
	}

	return int64(sum / count)
}
