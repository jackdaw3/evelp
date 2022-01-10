package esi

import (
	"encoding/json"
	"evelp/config/global"
	"evelp/model"
	"evelp/util"
	"fmt"
	"sort"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type OffersInit struct {
	offersMap map[int]*model.Offer
}

type offersWrapper struct {
	offers        *model.Offers
	corporationId int
}

func (o *OffersInit) Refresh() error {
	log.Infof("Start load offers from %s.", global.Conf.Data.RemoteDataAddress)
	o.offersMap = make(map[int]*model.Offer)
	if err := o.getOffersMap(); err != nil {
		return err
	}
	log.Info("Offers loaded.")

	log.Info("Save offers to DB.")
	var offers model.Offers
	for _, v := range o.offersMap {
		offers = append(offers, v)
	}
	sort.Sort(offers)
	if err := model.SaveOffers(&offers); err != nil {
		return err
	}
	log.Infof("%d offers have saved to DB.", offers.Len())

	return nil
}

func (o *OffersInit) getOffersMap() error {

	corporations, err := model.GetCorporations()
	if err != nil {
		return err
	}

	for _, corporation := range *corporations {
		wg.Add(1)
		go getOffers(corporation.CorporationId)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for offersWrapper := range channel {
		o.covertOffersWrapper(offersWrapper)
	}

	return nil
}

func (o *OffersInit) covertOffersWrapper(offersWrapper offersWrapper) {
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

func getOffers(corporationId int) {
	defer wg.Done()
	defer sem.Release(weigth)
	acquireSem(weigth)

	req := fmt.Sprintf("%s/loyalty/stores/%s/offers/?datasource=%s", global.Conf.Data.RemoteDataAddress, strconv.Itoa(corporationId), global.Conf.Data.RemoteDataSource)
	body, err := util.GetWithRetries(client, req)
	if err != nil {
		log.Errorf("Get corporation %d's failed: %s", corporationId, err.Error())
	}

	var offers model.Offers
	if err = json.Unmarshal(body, &offers); err != nil {
		log.Errorf("Unmarshal corporation %d's offers json failed: %s", corporationId, err.Error())
	}

	if offers.Len() == 0 {
		log.Warnf("Corporation %d has no offer.", corporationId)
		return
	}

	channel <- offersWrapper{&offers, corporationId}
}
