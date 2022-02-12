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

	"github.com/pkg/errors"
)

type itemHistroy struct {
	expirationTime time.Duration
}

func (i *itemHistroy) invoke() func() {
	return func() {
		products := make(map[int]struct{})

		offers, err := model.GetOffers()
		if err != nil {
			log.Error(err, "get offers failed")
			return
		}

		for _, offer := range *offers {
			if offer.IsBluePrint {
				bluePrint := model.GetBluePrint(offer.ItemId)
				if len(bluePrint.Products) == 0 {
					log.Error(errors.Errorf("offer %d's bluePrint %d have no product", offer.OfferId, bluePrint.BlueprintId))
					continue
				}
				product, err := model.GetItem(bluePrint.Products[0].ItemId)
				if err != nil {
					log.Errorf(err, "get item %d failed", bluePrint.Products[0].ItemId)
					continue
				}
				products[product.ItemId] = struct{}{}

			} else {
				products[offer.ItemId] = struct{}{}
			}
		}

		log.Infof("start load %d items history to redis", len(products))

		for p := range products {
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
				log.Errorf(err, "save orders to redis failed", key)
			}
		}

		log.Infof("Items histroy saved to reids")
	}
}
