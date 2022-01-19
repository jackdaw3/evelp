package dbdata

import (
	"evelp/initial/internal/api"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	wg     sync.WaitGroup
	client = &http.Client{}
)

func RemoteData() error {
	log.Info("start load remote data to DB")
	start := time.Now()

	offersInit := new(offersData)
	regionsInit := new(reginosData)
	starSystemsInit := new(starSystemsData)

	initializers := []api.Data{offersInit, regionsInit, starSystemsInit}
	for _, itinitializer := range initializers {
		if err := itinitializer.Refresh(); err != nil {
			return err
		}
	}

	elapsed := time.Since(start)
	log.Info("refresh remote data to DB cost: ", elapsed)
	return nil
}
