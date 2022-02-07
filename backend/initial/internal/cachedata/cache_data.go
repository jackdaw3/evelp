package cachedata

import (
	"evelp/model"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	orderExpireTime = 5 * time.Hour
)

func CacheData() error {
	log.Info("start refresh cache data")

	orders := make(map[string]*model.Orders)
	ordersData := new(ordersData)
	ordersData.orders = orders
	ordersData.expirationTime = orderExpireTime
	go func() {
		for {
			if err := ordersData.Refresh(); err != nil {
				log.Errorf("refresh orders to cache failed: %v", err)
			}
		}
	}()

	return nil
}
