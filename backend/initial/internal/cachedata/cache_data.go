package cachedata

import (
	"evelp/config/global"
	"evelp/initial/internal/api"
	"evelp/log"
	"evelp/model"
	"net/http"
	"time"
)

var client = &http.Client{}

const the_forge = 10000002

func CacheData() error {
	log.Info("start refresh cache data")

	orders := make(map[string]*model.Orders)
	ordersDataInit := new(ordersData)
	ordersDataInit.orders = orders
	ordersDataInit.expirationTime = global.Conf.Redis.ExpireTime.Order * time.Minute
	items, err := model.GetAllItems()
	if err != nil {
		return err
	}
	ordersDataInit.items = items

	itemHistroyDataInit := new(itemHistroy)
	itemHistroyDataInit.expirationTime = global.Conf.Redis.ExpireTime.History * time.Minute
	products, err := model.GetAllProducts()
	if err != nil {
		return err
	}
	itemHistroyDataInit.products = products

	initializers := []api.Data{ordersDataInit, itemHistroyDataInit}
	for _, itinitializer := range initializers {
		if err := itinitializer.Refresh(); err != nil {
			return err
		}
	}

	return nil
}
