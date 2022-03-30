package dto

type OrderDTO struct {
	OrderId      int
	ItemId       int
	ItemName     string
	SystemName   string
	VolumeRemain int64
	VolumeTotal  int64
	IsBuyOrder   bool
	Price        float64
	Expiration   string
	LastUpdated  string
}

type OrderDTOs []OrderDTO

func (o OrderDTOs) Len() int { return len(o) }

func (o OrderDTOs) Less(i, j int) bool { return o[i].Price < o[j].Price }

func (o OrderDTOs) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
