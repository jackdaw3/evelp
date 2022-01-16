package initial

import (
	"evelp/initial/internal/esi"
	"time"

	log "github.com/sirupsen/logrus"
)

func remoteData() error {
	log.Info("Start refresh remote esi data to DB.")
	start := time.Now()

	offersInit := new(esi.OffersInit)
	regionsInit := new(esi.ReginosInit)
	starSystemsInit := new(esi.StarSystemsInit)

	initializers := []esi.RemoteDataInit{offersInit, regionsInit, starSystemsInit}
	for _, itinitializer := range initializers {
		if err := itinitializer.Refresh(); err != nil {
			return err
		}
	}

	elapsed := time.Since(start)
	log.Info("Refresh remote data to DB cost: ", elapsed)
	return nil
}
