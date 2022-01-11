package esi

import (
	"encoding/json"
	"evelp/config/global"
	"evelp/model"
	"evelp/util/netutil"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
)

type ReginosInit struct {
	regions *model.Regions
}

func (r *ReginosInit) Refresh() error {
	log.Infof("Start load regions from %s.", global.Conf.Data.RemoteDataAddress)
	r.getAllRegions()
	sort.Sort(r.regions)

	for _, region := range *r.regions {
		exist, err := model.IsRegionExist(region.RegionId)
		if err != nil {
			log.Errorf("Check region %d exist failed %s.", region.RegionId, err)
		}

		if exist {
			continue
		}

		wg.Add(1)
		acquireSem(weigth)
		go getRegion(region)
	}
	wg.Wait()
	log.Info("Regions loaded and have saved to DB..")

	return nil
}

func (r *ReginosInit) getAllRegions() {
	req := fmt.Sprintf("%s/universe/regions/?datasource=%s", global.Conf.Data.RemoteDataAddress, global.Conf.Data.RemoteDataSource)

	body, err := netutil.GetWithRetries(client, req)
	if err != nil {
		log.Errorf("Get regions failed: %s", err.Error())
	}

	var idArray []int

	if err = json.Unmarshal(body, &idArray); err != nil {
		log.Errorf("Unmarshal regions json failed: %s", err.Error())
	}

	var regions model.Regions
	for _, id := range idArray {
		var region model.Region
		region.RegionId = id
		regions = append(regions, &region)
	}
	r.regions = &regions
}

func getRegion(region *model.Region) {
	defer wg.Done()
	defer sem.Release(weigth)

	for _, lang := range langs {
		req := fmt.Sprintf("%s/universe/regions/%d/?datasource=%s&language=%s", global.Conf.Data.RemoteDataAddress, region.RegionId, global.Conf.Data.RemoteDataSource, lang)

		body, err := netutil.GetWithRetries(client, req)
		if err != nil {
			log.Errorf("Get region %d failed: %s", region.RegionId, err.Error())
		}

		var resultMap map[string]interface{}

		if err = json.Unmarshal(body, &resultMap); err != nil {
			log.Errorf("Unmarshal region %d json failed: %s", region.RegionId, err.Error())
		}

		name, ok := resultMap["name"].(string)
		if !ok {
			log.Errorf("Region %d %v cast to string failed.", region.RegionId, resultMap["name"])
			continue
		}

		switch lang {
		case global.DE:
			region.Name.De = name
		case global.EN:
			region.Name.En = name
		case global.FR:
			region.Name.Fr = name
		case global.JA:
			region.Name.Ja = name
		case global.RU:
			region.Name.Ru = name
		case global.ZH:
			region.Name.Zh = name
		}
	}

	if err := model.SaveRegion(region); err != nil {
		log.Errorf("Region %d failed to save to DB.", region.RegionId)
	}
}
