package cachedata

import (
	"encoding/json"
	"evelp/config/global"
	"evelp/model"
	"evelp/util/cache"
	"evelp/util/net"
	"fmt"
	"io/ioutil"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	THE_FORGE = 10000002
	HISTROY   = "history"
)

type itemHistroy struct {
}

func (i *itemHistroy) invoke() func() {
	return func() {
		log.Infof("start load item history to redis")

		products := make(map[int]struct{})

		offers, err := model.GetOffers()
		if err != nil {
			log.Errorf("get offers failed: %+v", err)
			return
		}

		for _, offer := range *offers {
			if offer.IsBluePrint {
				bluePrint := model.GetBluePrint(offer.ItemId)
				if len(bluePrint.Products) == 0 {
					log.Errorf("offer %d's bluePrint %d have no product", offer.OfferId, bluePrint.BlueprintId)
				}
				product, err := model.GetItem(bluePrint.Products[0].ItemId)
				if err != nil {
					log.Errorf("get item %d failed: %+v", err)
				}
				products[product.ItemId] = struct{}{}

			} else {
				products[offer.ItemId] = struct{}{}
			}
		}

		for p := range products {
			req := fmt.Sprintf("%s/markets/%d/history/?datasource=%s&type_id=%d",
				global.Conf.Data.RemoteDataAddress,
				THE_FORGE,
				global.Conf.Data.RemoteDataSource,
				p,
			)

			resp, err := net.GetWithRetries(client, req)
			if err != nil {
				log.Errorf("get item %d histroy failed: %+v", p, err)

			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("read item %d histroybody failed: %+v", p, err)

			}

			var itemHistorys model.ItemHistorys
			if err = json.Unmarshal(body, &itemHistorys); err != nil {
				log.Errorf("unmarshal item %d histroy json failed: %+v", p, err)

			}

			key := cache.Key(ORDER, strconv.Itoa(THE_FORGE), strconv.Itoa(p))
			if err := cache.Set(key, itemHistorys, itemHistoryExpireTime); err != nil {
				log.Errorf("save orders to redis failed:%+v", key, err)
			}
		}

		log.Infof("Item histroy saved to reids")
	}
}
