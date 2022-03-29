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

	var localDataPath = global.Conf.Data.Local.Path
	typeIDsData := new(typeIdsData)
	typeIDsData.filePath = localDataPath + typeIdsFile

	factionsData := new(factionsData)
	factionsData.filePath = localDataPath + factionFile

	corporationsData := new(corporationsData)
	corporationsData.filePath = localDataPath + corporationsFile

	bluePrintsData := new(bluePrintsData)
	bluePrintsData.productFilePath = localDataPath + bluePrintProductsFile
	bluePrintsData.materialFilePath = localDataPath + bluePrintMaterialsFile

	localDataList := []api.Data{typeIDsData, factionsData, corporationsData, bluePrintsData}
	for _, localData := range localDataList {
		if err := localData.Refresh(); err != nil {
			return err
		}
	}

	elapsed := time.Since(start)
	log.Info("refresh static data to DB cost: ", elapsed)
	return nil
}
