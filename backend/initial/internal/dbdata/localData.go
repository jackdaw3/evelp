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
	log.Info("start to refresh static data to DB")
	start := time.Now()

	var localDataPath = global.Conf.Data.Local.Path
	typeIDData := new(typeIDData)
	typeIDData.filePath = localDataPath + typeIdsFile

	factionData := new(factionData)
	factionData.filePath = localDataPath + factionFile

	corporationData := new(corporationData)
	corporationData.filePath = localDataPath + corporationsFile

	bluePrintData := new(bluePrintData)
	bluePrintData.productFilePath = localDataPath + bluePrintProductsFile
	bluePrintData.materialFilePath = localDataPath + bluePrintMaterialsFile

	localDataList := []api.Data{typeIDData, factionData, corporationData, bluePrintData}
	for _, localData := range localDataList {
		if err := localData.Refresh(); err != nil {
			return err
		}
	}

	elapsed := time.Since(start)
	log.Info("refresh static data to DB cost: ", elapsed)
	return nil
}
