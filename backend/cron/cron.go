package cron

import (
	"evelp/service"

	log "github.com/sirupsen/logrus"
)

func Cron() error {
	log.Info("Start cron job.")

	orderService := new(service.OrderService)

	orderService.LoadOrders()()

	return nil
}
