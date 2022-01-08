package initial

import (
	"evelp/configs/global"
	"evelp/initial/internal/sde"
	"time"

	log "github.com/sirupsen/logrus"
)

func staticData() error {
	log.Info("Start refresh static data to DB.")
	start := time.Now()

	typeIDsInit := new(sde.TypeIDsInit)
	typeIDsInit.FilePath = global.Conf.Data.StaticDataPath + global.TypeIDsFileName

	factionsInit := new(sde.FactionsInit)
	factionsInit.FilePath = global.Conf.Data.StaticDataPath + global.FactionFileName

	corporationsInit := new(sde.CorporationsInit)
	corporationsInit.FilePath = global.Conf.Data.StaticDataPath + global.CorporationsFileName

	bluePrintsInit := new(sde.BluePrintsInit)
	bluePrintsInit.ProductFilePath = global.Conf.Data.StaticDataPath + global.BluePrintProducts
	bluePrintsInit.MaterialFilePath = global.Conf.Data.StaticDataPath + global.BluePrintMaterials

	initializers := []sde.StaticDataInit{typeIDsInit, factionsInit, corporationsInit, bluePrintsInit}
	for _, itinitializer := range initializers {
		if err := itinitializer.Refresh(); err != nil {
			return err
		}
	}

	elapsed := time.Since(start)
	log.Info("Refresh static data to DB cost: ", elapsed)
	return nil
}
