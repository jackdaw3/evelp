package service

import (
	"evelp/model"
	"evelp/util/cache"
	"strconv"

	"github.com/pkg/errors"
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
	if orders == nil {
		return 0, nil
	}

	price, err := orders.HighestBuyPrice(o.scope)
	if err != nil {
		return 0, errors.WithMessagef(err, "get orders %d highest buy price error", o.itemId)
	}

	return price, nil
}

func (o *OrderService) LowestSellPrice() (float64, error) {
	var orders *model.Orders
	orders, err := o.Orders()
	if err != nil {
		return 0, err
	}
	if orders == nil {
		return 0, nil
	}

	price, err := orders.LowestSellPrice(o.scope)
	if err != nil {
		return 0, errors.WithMessagef(err, "get orders %d lowest buy price error", o.itemId)
	}

	return price, nil
}

func (o *OrderService) Orders() (*model.Orders, error) {
	var orders model.Orders

	key := cache.Key(order, strconv.Itoa(o.regionId), strconv.Itoa(o.itemId))
	exist, err := cache.Exist(key)
	if err != nil {
		return nil, err
	}

	if exist {
		if err := cache.Get(key, &orders); err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &orders, nil
}
