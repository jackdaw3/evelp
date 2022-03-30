package cachedata

import (
	"encoding/json"
	"evelp/config/global"
	"evelp/log"
	"evelp/model"
	"evelp/util/cache"
	"evelp/util/net"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
	mu sync.Mutex
)

type orderData struct {
	orders         map[string]*model.Orders
	expirationTime time.Duration
	items          map[int]interface{}
}

func (o *orderData) Refresh() error {
	go func() {
		for {
			log.Debugf("start load orders to redis")

			log.Debugf("start load %d region's orders", the_forge)
			if err := o.ordersByRegion(the_forge); err != nil {
				log.Errorf(err, "load %d region orders failed", the_forge)
				o.clearMap()
				continue
			}

			log.Debugf("start load %d region's orders to redis", the_forge)
			for key, order := range o.orders {
				if err := cache.Set(key, *order, o.expirationTime); err != nil {
					log.Errorf(err, "save order %s to redis failed", key)
				}
			}

			o.clearMap()
			log.Debugf("orders saved to reids")
		}
	}()

	return nil
}

func (o *orderData) ordersByRegion(regionId int) error {
	pages, err := o.marketPages(regionId)
	if err != nil {
		return err
	}

	for i := 1; i <= pages; i++ {
		wg.Add(1)
		global.Ants.Submit(o.ordersByRegionPage(regionId, i))
	}

	wg.Wait()
	return nil
}

func (o *orderData) ordersByRegionPage(regionId int, page int) func() {
	return func() {
		defer wg.Done()

		req := fmt.Sprintf("%s/markets/%d/orders/?datasource=%s&order_type=all&page=%d",
			global.Conf.Data.Remote.Address,
			regionId,
			global.Conf.Data.Remote.DataSource,
			page,
		)

		resp, err := net.GetWithRetries(client, req)
		if err != nil {
			log.Errorf(err, "failed to get %d region %d page's orders", regionId, page)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf(err, "failed to read %d region %d page's orders body", regionId, page)
			return
		}

		var orders model.Orders
		if err = json.Unmarshal(body, &orders); err != nil {
			log.Errorf(err, "failed to unmarshal %d region %d page's orders json", regionId, page)
			return
		}

		for _, order := range orders {
			_, ok := o.items[order.ItemId]
			if !ok {
				continue
			}
			key := cache.Key("order", strconv.Itoa(regionId), strconv.Itoa(order.ItemId))
			o.syncPutToMap(key, &order)
		}
	}
}

func (o *orderData) marketPages(regionId int) (int, error) {
	req := fmt.Sprintf("%s/markets/%d/orders/?datasource=%s&order_type=all&page=%d",
		global.Conf.Data.Remote.Address,
		regionId,
		global.Conf.Data.Remote.DataSource,
		1,
	)

	resp, err := net.GetWithRetries(client, req)
	if err != nil {
		return 0, err
	}

	pages, err := strconv.Atoi(resp.Header.Get("x-pages"))
	if err != nil {
		return 0, err
	}

	return pages, nil

}

func (o *orderData) syncPutToMap(key string, order *model.Order) {
	defer mu.Unlock()
	mu.Lock()

	orders, ok := o.orders[key]
	order.LastUpdated = time.Now()
	if ok {
		val := append(*orders, *order)
		o.orders[key] = &val
	} else {
		val := model.Orders{*order}
		o.orders[key] = &val
	}
}

func (o *orderData) clearMap() {
	for k := range o.orders {
		delete(o.orders, k)
	}
}
