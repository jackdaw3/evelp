package cachedata

import (
	"evelp/model"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	client = &http.Client{}
)

const (
	orderExpireTime       = 5 * time.Hour
	itemHistoryExpireTime = 24 * time.Hour
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

	itemHistroyData := new(itemHistroy)
	cron := cron.New(cron.WithSeconds())
	if _, err := cron.AddFunc("@daily", itemHistroyData.invoke()); err != nil {
		return err
	}

	cron.Run()
	return nil
}
