package dbdata

import (
	"evelp/config/global"
	"evelp/initial/internal/api"
	"evelp/log"
	"time"
)

const (
	typeIdsFile            = "typeIDs.yaml"
	factionFile            = "factions.yaml"
	corporationsFile       = "npcCorporations.yaml"
	bluePrintProductsFile  = "industryActivityProducts.csv"
	bluePrintMaterialsFile = "industryActivityMaterials.csv"
)

func LocalData() error {
	log.Info("start refresh static data to DB")
	start := time.Now()

	var localDataPath = global.Conf.Data.LocalDataPath
	typeIDsInit := new(typeIdsData)
	typeIDsInit.filePath = localDataPath + typeIdsFile

	factionsInit := new(factionsData)
	factionsInit.filePath = localDataPath + factionFile

	corporationsInit := new(corporationsData)
	corporationsInit.filePath = localDataPath + corporationsFile

	bluePrintsInit := new(bluePrintsData)
	bluePrintsInit.productFilePath = localDataPath + bluePrintProductsFile
	bluePrintsInit.materialFilePath = localDataPath + bluePrintMaterialsFile

	initializers := []api.Data{typeIDsInit, factionsInit, corporationsInit, bluePrintsInit}
	for _, itinitializer := range initializers {
		if err := itinitializer.Refresh(); err != nil {
			return err
		}
	}

	elapsed := time.Since(start)
	log.Info("refresh static data to DB cost: ", elapsed)
	return nil
}
