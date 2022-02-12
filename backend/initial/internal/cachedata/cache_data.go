package cachedata

import (
	"evelp/config/global"
	"evelp/log"
	"evelp/model"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	client = &http.Client{}
)

const (
	orderExpireTime       = 5 * time.Hour
	itemHistoryExpireTime = 24 * time.Hour
	THE_FORGE             = 10000002
)

func CacheData() error {
	log.Info("start refresh cache data")

	orders := make(map[string]*model.Orders)
	ordersData := new(ordersData)
	ordersData.orders = orders
	ordersData.expirationTime = global.Conf.Redis.OrderExpireTime * time.Hour
	go func() {
		for {
			if err := ordersData.Refresh(); err != nil {
				log.Errorf(err, "refresh orders to cache failed")
			}
		}
	}()

	itemHistroyData := new(itemHistroy)
	itemHistroyData.expirationTime = global.Conf.Redis.HistoryExpireTime * time.Hour
	cron := cron.New(cron.WithSeconds())
	if _, err := cron.AddFunc("0 01 20 * * ?", itemHistroyData.invoke()); err != nil {
		return err
	}

	cron.Start()
	return nil
}
