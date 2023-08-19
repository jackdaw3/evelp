package dbdata

import (
	"encoding/json"
	"evelp/config/global"
	"evelp/log"
	"evelp/model"
	"evelp/util/net"
	"fmt"
	"io/ioutil"
	"sort"
	"sync"

	"github.com/pkg/errors"
)

type reginoData struct {
	regions *model.Regions
}

func (r *reginoData) Refresh() error {
	log.Infof("start to load regions from %s", global.Conf.Data.Remote.Address)
	r.getAllRegions()
	sort.Sort(r.regions)

	for _, region := range *r.regions {
		exist, err := region.IsExist()
		if err != nil {
			log.Errorf(err, "failed to check region %d exist", region.RegionId)
		}

		if exist {
			valid, err := region.IsVaild()
			if err != nil {
				log.Errorf(err, "failed to check region %d valid", region.RegionId)
			}

			if valid {
				continue
			}
		}

		wg.Add(1)
		if err := global.Ants.Submit(r.getRegion(region, &wg)); err != nil {
			return err
		}
	}

	wg.Wait()
	log.Info("regions loaded and saved to DB")

	return nil
}

func (r *reginoData) getAllRegions() {
	req := fmt.Sprintf("%s/universe/regions/?datasource=%s",
		global.Conf.Data.Remote.Address,
		global.Conf.Data.Remote.DataSource,
	)

	resp, err := net.Get(req)
	if err != nil {
		log.Errorf(err, "failed to get regions")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf(err, "failed to get regions' body")
	}

	var idArray []int

	if err = json.Unmarshal(body, &idArray); err != nil {
		log.Errorf(err, "failed to unmarshal regions json")
	}

	var regions model.Regions
	for _, id := range idArray {
		var region model.Region
		region.RegionId = id
		regions = append(regions, &region)
	}
	r.regions = &regions
}

func (r *reginoData) getRegion(region *model.Region, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()

		for _, lang := range global.Langs {
			req := fmt.Sprintf("%s/universe/regions/%d/?datasource=%s&language=%s",
				global.Conf.Data.Remote.Address,
				region.RegionId,
				global.Conf.Data.Remote.DataSource,
				lang,
			)

			resp, err := net.Get(req)
			if err != nil {
				log.Errorf(err, "failed to get region %d", region.RegionId)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf(err, "failed to get region %d's body", region.RegionId)
			}

			var resultMap map[string]interface{}

			if err = json.Unmarshal(body, &resultMap); err != nil {
				log.Errorf(err, "failed to unmarshal region %d json", region.RegionId)
			}

			name, ok := resultMap["name"].(string)
			if !ok {
				log.Error(errors.Errorf("failedt to cast region %d %v to string", region.RegionId, resultMap["name"]))
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
			log.Error(errors.Errorf("failed to save region %d to DB", region.RegionId))
		}
	}
}
