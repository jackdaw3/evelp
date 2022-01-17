package service

import (
	"encoding/json"
	"evelp/config/global"
	"evelp/model"
	"evelp/util/cache"
	"evelp/util/net"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	client = &http.Client{}
	wg     sync.WaitGroup
	mu     sync.Mutex
)

const (
	ORDER  = "order"
	XPAGES = "x-pages"
	FIRST  = 1
)

type OrderService struct {
	orders         map[string]*model.Orders
	expirationTime time.Duration
}

func NewOrderService(expirationTime time.Duration) *OrderService {
	orders := make(map[string]*model.Orders)
	return &OrderService{orders, expirationTime}
}

func (o *OrderService) LoadOrders() func() {
	return func() {
		regions, err := model.GetRegions()
		if err != nil {
			log.Errorf("Get regions failed.")
			return
		}

		for _, region := range *regions {
			log.Infof("Start load %d region's orders.", region.RegionId)
			if err := o.loadOrdersByRegion(region.RegionId); err != nil {
				log.Errorf("Load %d region orders failed:%v", region.RegionId, err)
				continue
			}

			log.Infof("Start save %d region's orders to redis.", region.RegionId)
			for key, order := range o.orders {
				if err := cache.Set(key, order, o.expirationTime); err != nil {
					log.Errorf("Save orders to redis failed:%v", key, err)
				}
			}

			o.clearMap()
		}
	}
}

func (o *OrderService) loadOrdersByRegion(regionId int) error {
	pages, err := o.getOrdersPage(regionId)
	if err != nil {
		return err
	}

	for i := 1; i <= pages; i++ {
		wg.Add(1)
		global.ANTS.Submit(o.loadOrdersByRegionPage(regionId, i))
	}

	wg.Wait()
	return nil
}

func (o *OrderService) loadOrdersByRegionPage(regionId int, page int) func() {
	return func() {
		defer wg.Done()

		req := fmt.Sprintf("%s/markets/%d/orders/?datasource=%s&order_type=all&page=%d",
			global.Conf.Data.RemoteDataAddress,
			regionId,
			global.Conf.Data.RemoteDataSource,
			page,
		)

		resp, err := net.GetWithRetries(client, req)
		if err != nil {
			log.Errorf("Get %d region %d page's orders failed: %v", regionId, page, err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Read %d region %d page's orders body failed: %v", regionId, page, err)
			return
		}

		var orders model.Orders
		if err = json.Unmarshal(body, &orders); err != nil {
			log.Errorf("Unmarshal %d region %d page's orders json failed: %v", regionId, page, err)
			return
		}

		for _, order := range orders {
			key := cache.Key(ORDER, strconv.Itoa(regionId), strconv.Itoa(order.ItemId))
			o.syncPutToMap(key, &order)
		}
	}
}

func (o *OrderService) getOrdersPage(regionId int) (int, error) {
	req := fmt.Sprintf("%s/markets/%d/orders/?datasource=%s&order_type=all&page=%d",
		global.Conf.Data.RemoteDataAddress,
		regionId,
		global.Conf.Data.RemoteDataSource,
		FIRST,
	)

	resp, err := net.GetWithRetries(client, req)
	if err != nil {
		return 0, err
	}

	pages, err := strconv.Atoi(resp.Header.Get(XPAGES))
	if err != nil {
		return 0, err
	}

	return pages, nil

}

func (o *OrderService) syncPutToMap(key string, order *model.Order) {
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

func (o *OrderService) clearMap() {
	for k := range o.orders {
		delete(o.orders, k)
	}
}
