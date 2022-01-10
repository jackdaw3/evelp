package esi

import (
	"context"
	"encoding/json"
	"evelp/config/global"
	"evelp/model"
	"evelp/util/netUtil"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
)

type ReginosInit struct {
	regions *model.Regions
}

var langs [6]string = [6]string{global.DE, global.EN, global.FR, global.JA, global.RU, global.ZH}

func (r *ReginosInit) Refresh() error {
	log.Infof("Start load regions from %s.", global.Conf.Data.RemoteDataAddress)
	r.getAllRegions()
	for _, region := range *r.regions {
		wg.Add(1)
		go getRegion(region)
	}
	wg.Wait()
	log.Info("Regions loaded.")

	log.Info("Save regions to DB.")
	sort.Sort(r.regions)
	if err := model.SaveRegions(r.regions); err != nil {
		return err
	}
	log.Infof("%d regions have saved to DB.", r.regions.Len())

	return nil
}

func (r *ReginosInit) getAllRegions() {
	req := fmt.Sprintf("%s/universe/regions/?datasource=%s", global.Conf.Data.RemoteDataAddress, global.Conf.Data.RemoteDataSource)

	body, err := netUtil.GetWithRetries(client, req)
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
	if err := sem.Acquire(context.Background(), weigth); err != nil {
		log.Errorf("Region %d get sem failed", region.RegionId, err.Error())
		return
	}

	for _, lang := range langs {
		req := fmt.Sprintf("%s/universe/regions/%d/?datasource=%s&language=%s", global.Conf.Data.RemoteDataAddress, region.RegionId, global.Conf.Data.RemoteDataSource, lang)

		body, err := netUtil.GetWithRetries(client, req)
		if err != nil {
			log.Errorf("Get region %d failed: %s", region.RegionId, err.Error())
		}

		var resultMap map[string]interface{}

		if err = json.Unmarshal(body, &resultMap); err != nil {
			log.Errorf("Unmarshal region %d json failed: %s", region.RegionId, err.Error())
		}

		name := resultMap["name"].(string)
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
}
