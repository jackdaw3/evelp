package cachedata

import (
	"evelp/config/global"
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
	ordersData := new(ordersData)
	ordersData.orders = orders
	ordersData.expirationTime = global.Conf.Redis.ExpireTime.Order * time.Minute
	items, err := model.GetAllItems()
	if err != nil {
		return err
	}
	ordersData.items = items

	go func() {
		for {
			if err := ordersData.Refresh(); err != nil {
				log.Errorf(err, "refresh orders to cache failed")
			}
		}
	}()

	itemHistroyData := new(itemHistroy)
	itemHistroyData.expirationTime = global.Conf.Redis.ExpireTime.History * time.Minute
	products, err := model.GetAllProducts()
	if err != nil {
		return err
	}
	itemHistroyData.products = products

	if err := itemHistroyData.Refresh(); err != nil {
		return err
	}

	return nil
}
