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

type starSystemsData struct {
	starSystems *model.StarSystems
}

func (s *starSystemsData) Refresh() error {
	log.Infof("start load starSystems from %s", global.Conf.Data.RemoteDataAddress)
	s.getAllStarSystems()
	sort.Sort(s.starSystems)

	for _, starSystem := range *s.starSystems {
		exist, err := starSystem.IsExist()
		if err != nil {
			log.Errorf("check starSystem %d exist failed: %v", starSystem.SystemId, err)
		}

		if exist {
			valid, err := starSystem.IsVaild()
			if err != nil {
				log.Errorf("check starSystem %d valid failed: %v", starSystem.SystemId, err)
			}

			if valid {
				continue
			}
		}

		wg.Add(1)
		global.ANTS.Submit(s.getStarSystem(starSystem, &wg))
	}

	wg.Wait()
	log.Info("starSystems have loaded and saved to DB")

	return nil
}

func (s *starSystemsData) getAllStarSystems() {
	req := fmt.Sprintf("%s/universe/systems/?datasource=%s",
		global.Conf.Data.RemoteDataAddress,
		global.Conf.Data.RemoteDataSource,
	)

	resp, err := net.GetWithRetries(client, req)
	if err != nil {
		log.Errorf("get starSystems failed: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("get starSystems' body failed: %v", err)
	}

	var idArray []int

	if err = json.Unmarshal(body, &idArray); err != nil {
		log.Errorf("unmarshal starSystems json failed: %v", err)
	}

	var starSystems model.StarSystems
	for _, id := range idArray {
		var starSystem model.StarSystem
		starSystem.SystemId = id
		starSystems = append(starSystems, &starSystem)
	}
	s.starSystems = &starSystems
}

func (s *starSystemsData) getStarSystem(starSystem *model.StarSystem, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()

		for _, lang := range global.LANGS {
			req := fmt.Sprintf("%s/universe/systems/%d/?datasource=%s&language=%s",
				global.Conf.Data.RemoteDataAddress,
				starSystem.SystemId,
				global.Conf.Data.RemoteDataSource,
				lang,
			)

			resp, err := net.GetWithRetries(client, req)
			if err != nil {
				log.Errorf("get starSystem %d failed: %v", starSystem.SystemId, err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("get starSystem %d's body failed: %v", starSystem.SystemId, err)
			}

			var resultMap map[string]interface{}

			if err = json.Unmarshal(body, &resultMap); err != nil {
				log.Errorf("unmarshal starSystem %d json failed: %v", starSystem.SystemId, err)
			}

			name, ok := resultMap["name"].(string)
			if !ok {
				log.Errorf("starSystem %d %v cast to string failed", starSystem.SystemId, resultMap["name"])
				continue
			}

			switch lang {
			case global.DE:
				starSystem.Name.De = name
			case global.EN:
				starSystem.Name.En = name
			case global.FR:
				starSystem.Name.Fr = name
			case global.JA:
				starSystem.Name.Ja = name
			case global.RU:
				starSystem.Name.Ru = name
			case global.ZH:
				starSystem.Name.Zh = name
			}
		}

		if err := model.SaveStarSystem(starSystem); err != nil {
			log.Errorf("starSystem %d save to DB failed: %v", starSystem.SystemId, err)
		}
	}
}
