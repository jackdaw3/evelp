package dto

type ItemHistoryDTO struct {
	ItemId     int
	Date       string
	Average    float64
	Average5d  float64
	Average20d float64
	Highest    float64
	Highest5d  float64
	Lowest     float64
	Lowest5d   float64
	OrderCount int64
	Volume     int64
}

type ItemHistoryDTOs []*ItemHistoryDTO

func (ihs *ItemHistoryDTOs) GenerateHistory() {
	for index, historyDTO := range *ihs {
		historyDTO.Average5d = ihs.average(index, 5)
		historyDTO.Average20d = ihs.average(index, 20)
		historyDTO.Highest5d = ihs.max(index, 5)
		historyDTO.Lowest5d = ihs.min(index, 5)
	}
}

func (ihs *ItemHistoryDTOs) average(index int, days int) float64 {
	var (
		sum   float64
		count int
		size  int
	)

	if index+1 < days {
		size = index + 1
	} else {
		size = days
	}

	for count < size {
		sum += (*ihs)[index-count].Average
		count++
	}

	if count == 0 {
		return 0
	}

	return sum / float64(count)
}

func (ihs *ItemHistoryDTOs) min(index int, days int) float64 {
	var (
		count int
		size  int
	)
	min := (*ihs)[index].Lowest

	if index+1 < days {
		size = index + 1
	} else {
		size = days
	}

	for count < size {
		tmp := (*ihs)[index-count].Lowest
		if min > tmp {
			min = tmp
		}
		count++
	}

	return min
}

func (ihs *ItemHistoryDTOs) max(index int, days int) float64 {
	var (
		max   float64
		count int
		size  int
	)

	if index+1 < days {
		size = index + 1
	} else {
		size = days
	}

	for count < size {
		tmp := (*ihs)[index-count].Highest
		if max < tmp {
			max = tmp
		}
		count++
	}

	return max
}
