package service

import (
	"evelp/dto"
	"evelp/log"
	"evelp/model"
	"evelp/util/cache"
	"sort"
	"strconv"

	"github.com/pkg/errors"
)

const order = "order"

type OrderService struct {
	itemId      int
	regionId    int
	isBluePrint bool
	scope       float64
}

func NewOrderService(itemId int, regionId int, isBluePrint bool, scope float64) *OrderService {
	return &OrderService{itemId, regionId, isBluePrint, scope}
}

func (o *OrderService) HighestBuyPrice() (float64, error) {
	var orders *model.Orders
	orders, err := o.ordersFromCache()
	if err != nil {
		return 0, err
	}

	price, err := orders.HighestBuyPrice(o.scope)
	if err != nil {
		return 0, errors.WithMessagef(err, "get orders %d highest buy price error", o.itemId)
	}

	return price, nil
}

func (o *OrderService) LowestSellPrice() (float64, error) {
	var orders *model.Orders
	orders, err := o.ordersFromCache()
	if err != nil {
		return 0, err
	}

	price, err := orders.LowestSellPrice(o.scope)
	if err != nil {
		return 0, errors.WithMessagef(err, "get orders %d lowest buy price error", o.itemId)
	}

	return price, nil
}

func (o *OrderService) Orders(isBuyOrder bool, lang string) (*dto.OrderDTOs, error) {
	orders, err := o.ordersFromCache()
	if err != nil {
		return nil, err
	}

	var orderDTOs dto.OrderDTOs

	item, err := model.GetItem(o.itemId)
	if err != nil {
		return nil, err
	}
	itemName := item.Name.Lang(lang)

	for _, order := range *orders {
		if !order.IsBuyOrder == isBuyOrder {
			continue
		}
		var orderDTO dto.OrderDTO
		orderDTO.OrderId = order.OrderId
		orderDTO.ItemId = order.ItemId
		orderDTO.ItemName = itemName
		orderDTO.Expiration = order.ExpirationToString()
		orderDTO.LastUpdated = order.LastUpdatedToString()
		orderDTO.IsBuyOrder = order.IsBuyOrder
		orderDTO.Price = order.Price
		orderDTO.VolumeRemain = order.VolumeRemain
		orderDTO.VolumeTotal = order.VolumeTotal

		system, err := model.GetStarSystem(order.SystemId)
		if err != nil {
			log.Errorf(err, "failed to get star system %v", order.SystemId)
			continue
		}
		orderDTO.SystemName = system.Name.Lang(lang)

		orderDTOs = append(orderDTOs, orderDTO)
	}

	if isBuyOrder {
		sort.Sort(sort.Reverse(orderDTOs))
	} else {
		sort.Sort(orderDTOs)
	}

	return &orderDTOs, nil
}

func (o *OrderService) ordersFromCache() (*model.Orders, error) {
	var orders model.Orders

	var key string
	if o.isBluePrint {
		bluePrint, err := model.GetBluePrint(o.itemId)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to get blue print %d", o.itemId)
		}
		if len(bluePrint.Products) == 0 {
			return nil, errors.Errorf("offer %d's bluePrint %d have no product", o.itemId, bluePrint.BlueprintId)
		}
		key = cache.Key(order, strconv.Itoa(o.regionId), strconv.Itoa(bluePrint.Products[0].ItemId))
	} else {
		key = cache.Key(order, strconv.Itoa(o.regionId), strconv.Itoa(o.itemId))
	}

	if err := cache.Get(key, &orders); err != nil {
		if errors.Cause(err).Error() == "redis: nil" {
			return nil, errors.New("no order found in the market")
		}
		return nil, err
	}

	return &orders, nil
}
