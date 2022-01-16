package esi

import (
	"encoding/json"
	"evelp/config/global"
	"evelp/model"
	"evelp/util/netUtil"
	"fmt"
	"sort"
	"sync"

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
		exist, err := region.IsExist()
		if err != nil {
			log.Errorf("Check region %d exist failed %s.", region.RegionId, err)
		}

		if exist {
			valid, err := region.IsVaild()
			if err != nil {
				log.Errorf("Check region %d valid failed %s.", region.RegionId, err)
			}

			if valid {
				continue
			}
		}

		wg.Add(1)
		if err := global.ANTS.Submit(r.getRegion(region, &wg)); err != nil {
			return err
		}
	}

	wg.Wait()
	log.Info("Regions have loaded and saved to DB.")

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

func (r *ReginosInit) getRegion(region *model.Region, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()

		for _, lang := range global.LANGS {
			req := fmt.Sprintf("%s/universe/regions/%d/?datasource=%s&language=%s", global.Conf.Data.RemoteDataAddress, region.RegionId, global.Conf.Data.RemoteDataSource, lang)

			body, err := netUtil.GetWithRetries(client, req)
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
}
