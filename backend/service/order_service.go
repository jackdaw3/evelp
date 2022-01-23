package service

import (
	"evelp/model"
	"evelp/util/cache"
	"fmt"
	"strconv"
)

const order = "order"

type OrderService struct {
	itemId   int
	regionId int
	scope    float64
}

func NewOrderService(itemId int, regionId int, scope float64) *OrderService {
	return &OrderService{itemId, regionId, scope}
}

func (o *OrderService) HighestBuyPrice() (float64, error) {
	var orders *model.Orders
	orders, err := o.Orders()
	if err != nil {
		return 0, err
	}

	price, err := orders.HighestBuyPrice(o.scope)
	if err != nil {
		return 0, fmt.Errorf("get orders %d highest buy price error: %v", o.itemId, err)
	}

	return price, nil
}

func (o *OrderService) LowestSellPrice() (float64, error) {
	var orders *model.Orders
	orders, err := o.Orders()
	if err != nil {
		return 0, err
	}

	price, err := orders.LowestSellPrice(o.scope)
	if err != nil {
		return 0, fmt.Errorf("get orders %d lowest sell price error: %v", o.itemId, err)
	}

	return price, nil
}

func (o *OrderService) Orders() (*model.Orders, error) {
	var orders model.Orders
	key := cache.Key(order, strconv.Itoa(o.regionId), strconv.Itoa(o.itemId))
	if err := cache.Get(key, &orders); err != nil {
		return nil, fmt.Errorf("get order %s cache error: %v", key, err)
	}
	return &orders, nil
}
