package initial

import (
	"evelp/service"
	"time"
)

const (
	ONE_SECOND = 1 * time.Second
	FIVE_HOUR  = 5 * time.Hour
)

func schedule() error {
	orderService := service.NewOrderService(FIVE_HOUR)
	orderSchedule := service.NewScheduleService(orderService.LoadOrders(), ONE_SECOND)

	schedules := service.NewScheduleServices(*orderSchedule)
	if err := schedules.Start(); err != nil {
		return err
	}

	return nil
}
