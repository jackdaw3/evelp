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
)

type itemHistroy struct {
	expirationTime time.Duration
	products       map[int]interface{}
}

func (i *itemHistroy) invoke() func() {
	return func() {
		log.Debug("start load items history to redis")

		for p := range i.products {
			req := fmt.Sprintf("%s/markets/%d/history/?datasource=%s&type_id=%d",
				global.Conf.Data.Remote.Address,
				the_forge,
				global.Conf.Data.Remote.DataSource,
				p,
			)

			resp, err := net.GetWithRetries(client, req)
			if err != nil {
				log.Errorf(err, "get item %d histroy failed", p)
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf(err, "read item %d histroybody failed", p)
				continue
			}

			var itemHistorys model.ItemHistorys
			if err = json.Unmarshal(body, &itemHistorys); err != nil {
				log.Errorf(err, "unmarshal item %d histroy json failed", p)
				continue
			}

			for _, itemitemHistory := range itemHistorys {
				itemitemHistory.ItemId = p
			}

			key := cache.Key("history", strconv.Itoa(the_forge), strconv.Itoa(p))
			if err := cache.Set(key, itemHistorys, i.expirationTime); err != nil {
				log.Errorf(err, "save orders %v to redis failed", key)
			}
		}

		log.Debugf("Items histroy saved to reids")
	}
}
