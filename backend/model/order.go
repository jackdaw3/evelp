package model

import (
	"errors"
	"math"
	"sort"
	"time"
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

func (orders Orders) Len() int { return len(orders) }

func (orders Orders) Less(i, j int) bool { return orders[i].Price < orders[j].Price }

func (orders Orders) Swap(i, j int) { orders[i], orders[j] = orders[j], orders[i] }

func (o *Orders) GetHighestBuyPrice(scope float64) (float64, error) {
	if err := o.isValid(); err != nil {
		return 0, err
	}

	var buyOrders Orders
	for _, order := range *o {
		if order.IsBuyOrder {
			buyOrders = append(buyOrders, order)
		}
	}
	sort.Sort(sort.Reverse(buyOrders))

	return buyOrders.getOrdersPrice(scope)
}

func (o *Orders) GetLowestSellPrice(scope float64) (float64, error) {
	if err := o.isValid(); err != nil {
		return 0, err
	}

	var sellOrders Orders
	for _, order := range *o {
		if !order.IsBuyOrder {
			sellOrders = append(sellOrders, order)
		}
	}
	sort.Sort(sellOrders)

	return sellOrders.getOrdersPrice(scope)
}

func (o *Orders) getOrdersPrice(scope float64) (float64, error) {
	size, err := o.getOrdersScopeSize(scope)
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

func (o *Orders) getOrdersScopeSize(scope float64) (int64, error) {
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
