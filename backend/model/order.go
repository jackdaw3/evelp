package model

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type Order struct {
	OrderId      int
	ItemId       int
	Issued       time.Time
	Duration     int
	SystemId     int
	Price        float64
	VolumeRemain int64
	VolumeTotal  int64
	IsBuyOrder   bool
}

type Orders []Order

func (orders Orders) Len() int { return len(orders) }

func (orders Orders) Less(i, j int) bool { return orders[i].Price < orders[j].Price }

func (orders Orders) Swap(i, j int) { orders[i], orders[j] = orders[j], orders[i] }

func (o *Orders) GetHighestBuyPrice(scope float64) (float64, error) {
	var buyOrders Orders
	for _, order := range *o {
		if order.IsBuyOrder {
			buyOrders = append(buyOrders, order)
		}
	}
	sort.Sort(sort.Reverse(buyOrders))

	return getOrdersPrice(scope, &buyOrders)
}

func (o *Orders) GetLowestSellPrice(scope float64) (float64, error) {
	var sellOrders Orders
	for _, order := range *o {
		if !order.IsBuyOrder {
			sellOrders = append(sellOrders, order)
		}
	}
	sort.Sort(sellOrders)

	return getOrdersPrice(scope, &sellOrders)
}

func getOrdersPrice(scope float64, orders *Orders) (float64, error) {
	size, err := getScopeSize(scope, orders)
	if err != nil {
		return 0, err
	}

	var sum float64
	var count int64
	for i := 0; i < len(*orders) && count < size; i++ {
		volume := (*orders)[i].VolumeRemain
		price := (*orders)[i].Price

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

func getScopeSize(scope float64, orders *Orders) (int64, error) {
	var size int64
	for _, order := range *orders {
		size += order.VolumeRemain
	}

	if size <= 0 || scope <= 0 {
		return 0, fmt.Errorf("GetScopeSize incorrect paratermer with size:%d and scope:%f", size, scope)
	}

	result := int64(math.Ceil(float64(size) * scope))
	if result == 0 {
		return 1, nil
	} else {
		return result, nil
	}
}
