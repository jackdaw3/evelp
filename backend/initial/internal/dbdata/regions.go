package dbdata

import (
	"encoding/json"
	"evelp/config/global"
	"evelp/model"
	"evelp/util/net"
	"fmt"
	"io/ioutil"
	"sort"
	"sync"

	log "github.com/sirupsen/logrus"
)

type reginosData struct {
	regions *model.Regions
}

func (r *reginosData) Refresh() error {
	log.Infof("start load regions from %s", global.Conf.Data.RemoteDataAddress)
	r.getAllRegions()
	sort.Sort(r.regions)

	for _, region := range *r.regions {
		exist, err := region.IsExist()
		if err != nil {
			log.Errorf("check region %d exist failed: %v", region.RegionId, err)
		}

		if exist {
			valid, err := region.IsVaild()
			if err != nil {
				log.Errorf("check region %d valid failed: %v", region.RegionId, err)
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
	log.Info("regions have loaded and saved to DB")

	return nil
}

func (r *reginosData) getAllRegions() {
	req := fmt.Sprintf("%s/universe/regions/?datasource=%s",
		global.Conf.Data.RemoteDataAddress,
		global.Conf.Data.RemoteDataSource,
	)

	resp, err := net.GetWithRetries(client, req)
	if err != nil {
		log.Errorf("get regions failed: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("get regions'body failed: %v", err)
	}

	var idArray []int

	if err = json.Unmarshal(body, &idArray); err != nil {
		log.Errorf("unmarshal regions json failed: %v", err)
	}

	var regions model.Regions
	for _, id := range idArray {
		var region model.Region
		region.RegionId = id
		regions = append(regions, &region)
	}
	r.regions = &regions
}

func (r *reginosData) getRegion(region *model.Region, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()

		for _, lang := range global.LANGS {
			req := fmt.Sprintf("%s/universe/regions/%d/?datasource=%s&language=%s",
				global.Conf.Data.RemoteDataAddress,
				region.RegionId,
				global.Conf.Data.RemoteDataSource,
				lang,
			)

			resp, err := net.GetWithRetries(client, req)
			if err != nil {
				log.Errorf("get region %d failed: %v", region.RegionId, err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("get region %d's body failed: %v", region.RegionId, err)
			}

			var resultMap map[string]interface{}

			if err = json.Unmarshal(body, &resultMap); err != nil {
				log.Errorf("unmarshal region %d json failed: %v", region.RegionId, err)
			}

			name, ok := resultMap["name"].(string)
			if !ok {
				log.Errorf("region %d %v cast to string failed", region.RegionId, resultMap["name"])
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
			log.Errorf("region %d failed to save to DB", region.RegionId)
		}
	}
}
