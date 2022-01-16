package service

import "evelp/model"

type OrderService struct {
	offersMap map[int]*model.Offers
}

func (o *OrderService) LoadOrders() error {

	return nil
}
