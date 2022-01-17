package esi

import (
	"encoding/json"
	"evelp/config/global"
	"evelp/model"
	"evelp/util/net"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

type OffersInit struct {
	offersMap map[int]*model.Offer
}

type offersWrapper struct {
	offers        *model.Offers
	corporationId int
}

var mu sync.Mutex

func (o *OffersInit) Refresh() error {
	log.Infof("Start load offers from %s.", global.Conf.Data.RemoteDataAddress)
	o.offersMap = make(map[int]*model.Offer)
	if err := o.getOffersMap(); err != nil {
		return err
	}
	log.Info("Offers have loaded.")

	log.Info("Start save offers to DB.")
	var offers model.Offers
	for _, v := range o.offersMap {
		offers = append(offers, v)
	}
	sort.Sort(offers)
	if err := model.SaveOffers(&offers); err != nil {
		return err
	}
	log.Infof("%d offers have saved or updated to DB.", offers.Len())

	return nil
}

func (o *OffersInit) getOffersMap() error {
	corporations, err := model.GetCorporations()
	if err != nil {
		return err
	}

	for _, corporation := range *corporations {
		wg.Add(1)
		if err := global.ANTS.Submit(o.getOffers(corporation.CorporationId, &wg)); err != nil {
			return err
		}
	}

	wg.Wait()
	return nil
}

func (o *OffersInit) covertOffersWrapper(offersWrapper offersWrapper) {
	defer mu.Unlock()
	mu.Lock()

	for _, offer := range *offersWrapper.offers {
		if value, ok := o.offersMap[offer.OfferId]; !ok {
			offer.CorporationIDs = append(offer.CorporationIDs, offersWrapper.corporationId)
			if bluePrint := model.GetBluePrint(offer.ItemId); !bluePrint.Empty() {
				offer.IsBluePrint = true
			}
			o.offersMap[offer.OfferId] = offer
		} else {
			hasCurrentCorporationId := false
			for _, v := range value.CorporationIDs {
				if v == offersWrapper.corporationId {
					hasCurrentCorporationId = true
				}
			}
			if !hasCurrentCorporationId {
				value.CorporationIDs = append(value.CorporationIDs, offersWrapper.corporationId)
			}
		}
	}
}

func (o *OffersInit) getOffers(corporationId int, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()

		req := fmt.Sprintf("%s/loyalty/stores/%s/offers/?datasource=%s",
			global.Conf.Data.RemoteDataAddress,
			strconv.Itoa(corporationId),
			global.Conf.Data.RemoteDataSource,
		)

		resp, err := net.GetWithRetries(client, req)
		if err != nil {
			log.Errorf("Get corporation %d's failed: %v", corporationId, err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Get corporation %d's body failed: %v", corporationId, err)
		}

		var offers model.Offers
		if err = json.Unmarshal(body, &offers); err != nil {
			log.Errorf("Unmarshal corporation %d's offers json failed: %v", corporationId, err)
		}

		if offers.Len() == 0 {
			log.Debugf("Corporation %d has no offer.", corporationId)
			return
		}

		o.covertOffersWrapper(offersWrapper{&offers, corporationId})
	}
}
