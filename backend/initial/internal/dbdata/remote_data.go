package dbdata

import (
	"evelp/initial/internal/api"
	"evelp/log"
	"net/http"
	"sync"
	"time"
)

var (
	wg     sync.WaitGroup
	client = &http.Client{}
)

func RemoteData() error {
	log.Info("start load remote data to DB")
	start := time.Now()

	offersData := new(offersData)
	regionsData := new(reginosData)
	starSystemsData := new(starSystemsData)

	remoteDataList := []api.Data{offersData, regionsData, starSystemsData}
	for _, remoteData := range remoteDataList {
		if err := remoteData.Refresh(); err != nil {
			return err
		}
	}

	elapsed := time.Since(start)
	log.Info("refresh remote data to DB cost: ", elapsed)
	return nil
}
