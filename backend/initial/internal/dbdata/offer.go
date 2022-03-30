package dbdata

import (
	"encoding/json"
	"evelp/config/global"
	"evelp/log"
	"evelp/model"
	"evelp/util/net"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"sync"
)

var mu sync.Mutex

type offerData struct {
	offersMap map[int]*model.Offer
}

type offersWrapper struct {
	offers        *model.Offers
	corporationId int
}

func (o *offerData) Refresh() error {
	log.Infof("start to load offers from %s", global.Conf.Data.Remote.Address)
	o.offersMap = make(map[int]*model.Offer)
	if err := o.getOffersMap(); err != nil {
		return err
	}
	log.Info("offers loaded")

	log.Info("start to save offers to DB")
	var offers model.Offers
	for _, v := range o.offersMap {
		offers = append(offers, v)
	}
	sort.Sort(offers)
	if err := model.SaveOffers(&offers); err != nil {
		return err
	}
	log.Infof("%d offers saved or updated to DB", offers.Len())

	return nil
}

func (o *offerData) getOffersMap() error {
	corporations, err := model.GetCorporations()
	if err != nil {
		return err
	}

	for _, corporation := range *corporations {
		wg.Add(1)
		if err := global.Ants.Submit(o.getOffers(corporation.CorporationId, &wg)); err != nil {
			return err
		}
	}

	wg.Wait()
	return nil
}

func (o *offerData) covertOffersWrapper(offersWrapper offersWrapper) {
	defer mu.Unlock()
	mu.Lock()

	for _, offer := range *offersWrapper.offers {
		if value, ok := o.offersMap[offer.OfferId]; !ok {
			offer.CorporationIds = append(offer.CorporationIds, offersWrapper.corporationId)
			bluePrint, err := model.GetBluePrint(offer.ItemId)
			if err != nil {
				log.Errorf(err, "get blue print %d failed", offer.ItemId)
				continue
			}
			if !bluePrint.Empty() {
				offer.IsBluePrint = true
			}
			o.offersMap[offer.OfferId] = offer
		} else {
			hasCurrentCorporationId := false
			for _, v := range value.CorporationIds {
				if v == offersWrapper.corporationId {
					hasCurrentCorporationId = true
				}
			}
			if !hasCurrentCorporationId {
				value.CorporationIds = append(value.CorporationIds, offersWrapper.corporationId)
			}
		}
	}
}

func (o *offerData) getOffers(corporationId int, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()

		req := fmt.Sprintf("%s/loyalty/stores/%s/offers/?datasource=%s",
			global.Conf.Data.Remote.Address,
			strconv.Itoa(corporationId),
			global.Conf.Data.Remote.DataSource,
		)

		resp, err := net.GetWithRetries(client, req)
		if err != nil {
			log.Errorf(err, "get corporation %d's failed", corporationId)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf(err, "get corporation %d's body failed", corporationId)
		}

		var offers model.Offers
		if err = json.Unmarshal(body, &offers); err != nil {
			log.Errorf(err, "unmarshal corporation %d's offers json failed", corporationId)
		}

		if offers.Len() == 0 {
			log.Debugf("corporation %d has no offer", corporationId)
			return
		}

		o.covertOffersWrapper(offersWrapper{&offers, corporationId})
	}
}
