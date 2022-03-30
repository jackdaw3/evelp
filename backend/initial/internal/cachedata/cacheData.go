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
	log.Info("start to refresh cache data")

	orders := make(map[string]*model.Orders)
	ordersData := new(orderData)
	ordersData.orders = orders
	ordersData.expirationTime = global.Conf.Redis.ExpireTime.Order * time.Minute
	items, err := model.GetAllItems()
	if err != nil {
		return err
	}
	ordersData.items = items

	itemHistroyData := new(itemHistroy)
	itemHistroyData.expirationTime = global.Conf.Redis.ExpireTime.History * time.Minute
	products, err := model.GetAllProducts()
	if err != nil {
		return err
	}
	itemHistroyData.products = products

	cachedataList := []api.Data{ordersData, itemHistroyData}
	for _, cachedata := range cachedataList {
		if err := cachedata.Refresh(); err != nil {
			return err
		}
	}

	return nil
}
