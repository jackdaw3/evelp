package dbdata

import (
	"evelp/initial/internal/api"
	"evelp/log"
	"sync"
	"time"
)

var wg sync.WaitGroup

func RemoteData() error {
	log.Info("start to load remote data to DB")
	start := time.Now()

	offerData := new(offerData)
	regionData := new(reginoData)
	starSystemData := new(starSystemData)

	remoteDataList := []api.Data{offerData, regionData, starSystemData}
	for _, remoteData := range remoteDataList {
		if err := remoteData.Refresh(); err != nil {
			return err
		}
	}

	elapsed := time.Since(start)
	log.Info("refresh remote data to DB cost: ", elapsed)
	return nil
}
