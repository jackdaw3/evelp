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
	"time"

	"github.com/robfig/cron/v3"
)

type itemHistroy struct {
	expirationTime time.Duration
	products       map[int]interface{}
}

func (i *itemHistroy) Refresh() error {
	cron := cron.New(cron.WithSeconds())
	if _, err := cron.AddFunc("@daily", i.invoke()); err != nil {
		return err
	}
	cron.Start()

	return nil
}

func (i *itemHistroy) invoke() func() {
	return func() {
		log.Debug("start to load items history to redis")

		for p := range i.products {
			req := fmt.Sprintf("%s/markets/%d/history/?datasource=%s&type_id=%d",
				global.Conf.Data.Remote.Address,
				the_forge,
				global.Conf.Data.Remote.DataSource,
				p,
			)

			resp, err := net.GetWithRetries(req)
			if err != nil {
				log.Errorf(err, "failed to get item %d histroy", p)
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf(err, "failed to read item %d histroybody ", p)
				continue
			}

			var itemHistorys model.ItemHistorys
			if err = json.Unmarshal(body, &itemHistorys); err != nil {
				log.Errorf(err, "failed to unmarshal item %d histroy json", p)
				continue
			}

			for _, itemitemHistory := range itemHistorys {
				itemitemHistory.ItemId = p
			}

			key := cache.Key("history", strconv.Itoa(the_forge), strconv.Itoa(p))
			if err := cache.Set(key, itemHistorys, i.expirationTime); err != nil {
				log.Errorf(err, "failed to save orders %v to redis", key)
			}
		}

		log.Debugf("Items histroy saved to reids")
	}
}
