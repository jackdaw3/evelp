package initial

import (
	"evelp/initial/internal/esi"
	"time"

	log "github.com/sirupsen/logrus"
)

func remoteData() error {
	log.Info("Start refresh remote esi data to DB.")
	start := time.Now()

	regionsInit := new(esi.ReginosInit)
	offersInit := new(esi.OffersInit)

	initializers := []esi.RemoteDataInit{regionsInit, offersInit}
	for _, itinitializer := range initializers {
		if err := itinitializer.Refresh(); err != nil {
			return err
		}
	}

	elapsed := time.Since(start)
	log.Info("Refresh static data to DB cost: ", elapsed)
	return nil
}
