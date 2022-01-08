package esi

import (
	"context"
	"encoding/json"
	"evelp/configs/global"
	"evelp/model"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

type OffersInit struct {
	offersMap map[int]*model.Offer
}

type offersWrapper struct {
	offers        *model.Offers
	corporationId int
}

var (
	buffCount int   = 5
	limit     int64 = 3
	weigth    int64 = 1
)

var sem = semaphore.NewWeighted(limit)
var channel = make(chan offersWrapper, buffCount)
var wg sync.WaitGroup
var client = &http.Client{}

func (offersInit *OffersInit) Refrsh() error {
	log.Infof("Start load offers from %s.", global.Conf.Data.RemoteDataAddress)
	offersInit.offersMap = make(map[int]*model.Offer)
	if err := offersInit.getOffersMap(); err != nil {
		return err
	}
	log.Info("Offers loaded.")

	log.Info("Save offers to DB.")
	var offers model.Offers
	for _, v := range offersInit.offersMap {
		offers = append(offers, v)
	}
	sort.Sort(offers)
	if err := model.SaveOffers(&offers); err != nil {
		return err
	}
	log.Infof("%d offers have saved to DB.", offers.Len())

	return nil
}

func (offersInit *OffersInit) getOffersMap() error {

	corporations, err := model.GetCorporations()
	if err != nil {
		return err
	}

	for _, corporation := range *corporations {
		wg.Add(1)
		go get(corporation.CorporationId)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for offersWrapper := range channel {
		offersInit.covertOffersWrapper(offersWrapper)
	}

	return nil
}

func (offersInit *OffersInit) covertOffersWrapper(offersWrapper offersWrapper) {
	for _, offer := range *offersWrapper.offers {
		if value, ok := offersInit.offersMap[offer.OfferId]; !ok {
			offer.CorporationIDs = append(offer.CorporationIDs, offersWrapper.corporationId)
			if bluePrint := model.GetBluePrint(offer.ItemId); !bluePrint.Empty() {
				offer.IsBluePrint = true
			}
			offersInit.offersMap[offer.OfferId] = offer
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

func get(corporationId int) {
	defer wg.Done()
	defer sem.Release(weigth)
	if err := sem.Acquire(context.Background(), weigth); err != nil {
		log.Errorf("Corporation %d get sem failed", corporationId, err.Error())
		return
	}

	req := fmt.Sprintf("%s/loyalty/stores/%s/offers/?datasource=%s", global.Conf.Data.RemoteDataAddress, strconv.Itoa(corporationId), global.Conf.Data.RemoteDataSource)
	resp, err := client.Get(req)
	if err != nil {
		log.Errorf("Get corporation %d's offers failed: %s", corporationId, err.Error())
		log.Warnf("Sleep 1 second and retry get corporation %d's offer.", corporationId)
		time.Sleep(1 * time.Second)
		resp, err = client.Get(req)
		if err != nil {
			log.Errorf("Retry get corporation %d's offers  failed: %s", corporationId, err.Error())
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Read corporation %d's offers body failed: %s", corporationId, err.Error())
	}

	var offers model.Offers
	if err = json.Unmarshal([]byte(body), &offers); err != nil {
		log.Errorf("Unmarshal corporation %d's offers json failed: %s", corporationId, err.Error())
	}

	if offers.Len() == 0 {
		log.Warnf("Corporation %d has no offer.", corporationId)
		return
	}

	channel <- offersWrapper{&offers, corporationId}
}
