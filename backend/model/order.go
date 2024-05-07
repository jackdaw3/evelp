package model

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/pkg/errors"
)

type Order struct {
	OrderId      int       `json:"order_id"`
	ItemId       int       `json:"type_id"`
	Issued       time.Time `json:"issued"`
	Duration     int       `json:"duration"`
	SystemId     int       `json:"system_id"`
	Price        float64   `json:"price"`
	VolumeRemain int64     `json:"volume_remain"`
	VolumeTotal  int64     `json:"volume_total"`
	IsBuyOrder   bool      `json:"is_buy_order"`
	LastUpdated  time.Time `json:"last_updated"`
}

type Orders []Order

func (o Orders) Len() int { return len(o) }

func (o Orders) Less(i, j int) bool { return o[i].Price < o[j].Price }

func (o Orders) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

func (o *Orders) HighestBuyPrice(scope float64) (float64, error) {
	if err := o.isValid(); err != nil {
		return 0, err
	}

	var buyOrders Orders
	for _, order := range *o {
		if order.IsBuyOrder {
			buyOrders = append(buyOrders, order)
		}
	}
	if len(buyOrders) == 0 {
		return 0, errors.New("no buy order found in the market")
	}
	sort.Sort(sort.Reverse(buyOrders))
	buyOrders = filterBuyOrders(buyOrders)

	return buyOrders.ordersPrice(scope)
}

func (o *Orders) LowestSellPrice(scope float64) (float64, error) {
	if err := o.isValid(); err != nil {
		return 0, err
	}

	var sellOrders Orders
	for _, order := range *o {
		if !order.IsBuyOrder {
			sellOrders = append(sellOrders, order)
		}
	}
	if len(sellOrders) == 0 {
		return 0, errors.New("no sell order found in the market")
	}
	sort.Sort(sellOrders)
	return sellOrders.ordersPrice(scope)
}

func (o *Order) LastUpdatedToString() string {
	var (
		value  string
		hour   int64
		minute int64
		second int64
	)

	duration := int64(time.Since(o.LastUpdated).Seconds())
	hour = duration / 3600
	minute = (duration % 3600) / 60
	second = duration % 60

	if hour != 0 {
		value = fmt.Sprintf("%dh ago", hour)
	} else if minute != 0 {
		value = fmt.Sprintf("%dm ago", minute)
	} else if second != 0 {
		value = fmt.Sprintf("%ds ago", second)
	}

	return value
}

func (o *Order) ExpirationToString() string {
	var (
		value  string
		day    int64
		hour   int64
		minute int64
		second int64
	)

	expirationTime := o.Issued.Add(time.Hour * 24 * time.Duration(o.Duration))
	if expirationTime.Before(time.Now()) {
		return "expired"
	}
	duration := int64(time.Until(expirationTime).Seconds())

	day = duration / (3600 * 24)
	hour = (duration % (3600 * 24)) / 3600
	minute = (duration % 3600) / 60
	second = duration % 60

	if day != 0 {
		value = fmt.Sprintf("%dd %dh %dm %ds later", day, hour, minute, second)
	} else if hour != 0 {
		value = fmt.Sprintf("%dh %dm %ds later", hour, minute, second)
	} else if minute != 0 {
		value = fmt.Sprintf("%dm %ds later", minute, second)
	} else if second != 0 {
		value = fmt.Sprintf("%ds later", second)
	}

	return value
}

func (o *Orders) ordersPrice(scope float64) (float64, error) {
	size, err := o.ordersScopeSize(scope)
	if err != nil {
		return 0, err
	}
	if size == 0 {
		return 0, nil
	}

	var sum float64
	var count int64
	for i := 0; i < len(*o) && count < size; i++ {
		volume := (*o)[i].VolumeRemain
		price := (*o)[i].Price

		if count+volume > size {
			volume = size - count
			count = size
		} else {
			count += volume
		}

		sum += float64(volume) * price
	}

	return sum / float64(size), nil
}

func (o *Orders) ordersScopeSize(scope float64) (int64, error) {
	var size int64
	for _, order := range *o {
		size += order.VolumeRemain
	}

	if size == 0 || scope == 0 {
		return 0, nil
	}

	result := int64(math.Ceil(float64(size) * scope))
	if result == 0 {
		return 1, nil
	} else {
		return result, nil
	}
}

func (o *Orders) isValid() error {
	var itemId int = (*o)[0].ItemId
	for _, order := range *o {
		if order.ItemId != itemId {
			return errors.New("orders have multiple itemIds")
		}
	}
	return nil
}

func filterBuyOrders(orders Orders) Orders {
	if len(orders) == 0 {
		return nil
	}

	var count int
	var currentVolume int64
	var maxPrice float64 = orders[0].Price

	for i, order := range orders {
		if order.IsBuyOrder {
			if maxPrice >= 1000*order.Price || (maxPrice >= 100*order.Price && order.VolumeRemain >= 20*currentVolume) {
				count = i
				break
			}
			currentVolume += order.VolumeRemain
		}
	}

	if count != 0 {
		filtedOrders := orders[:count]
		return filtedOrders
	}

	return orders
}
