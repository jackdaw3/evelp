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

type starSystemsData struct {
	starSystems *model.StarSystems
}

func (s *starSystemsData) Refresh() error {
	log.Infof("start load starSystems from %s", global.Conf.Data.Remote.Address)
	s.getAllStarSystems()
	sort.Sort(s.starSystems)

	for _, starSystem := range *s.starSystems {
		exist, err := starSystem.IsExist()
		if err != nil {
			log.Errorf(err, "check starSystem %d exist failed", starSystem.SystemId)
		}

		if exist {
			valid, err := starSystem.IsVaild()
			if err != nil {
				log.Errorf(err, "check starSystem %d valid failed", starSystem.SystemId)
			}

			if valid {
				continue
			}
		}

		wg.Add(1)
		global.Ants.Submit(s.getStarSystem(starSystem, &wg))
	}

	wg.Wait()
	log.Info("starSystems have loaded and saved to DB")

	return nil
}

func (s *starSystemsData) getAllStarSystems() {
	req := fmt.Sprintf("%s/universe/systems/?datasource=%s",
		global.Conf.Data.Remote.Address,
		global.Conf.Data.Remote.DataSource,
	)

	resp, err := net.GetWithRetries(client, req)
	if err != nil {
		log.Errorf(err, "get starSystems failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf(err, "get starSystems' body failed")
	}

	var idArray []int

	if err = json.Unmarshal(body, &idArray); err != nil {
		log.Errorf(err, "unmarshal starSystems json failed")
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

		for _, lang := range global.Langs {
			req := fmt.Sprintf("%s/universe/systems/%d/?datasource=%s&language=%s",
				global.Conf.Data.Remote.Address,
				starSystem.SystemId,
				global.Conf.Data.Remote.DataSource,
				lang,
			)

			resp, err := net.GetWithRetries(client, req)
			if err != nil {
				log.Errorf(err, "get starSystem %d failed", starSystem.SystemId)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf(err, "get starSystem %d's body failed", starSystem.SystemId)
			}

			var resultMap map[string]interface{}

			if err = json.Unmarshal(body, &resultMap); err != nil {
				log.Errorf(err, "unmarshal starSystem %d json failed", starSystem.SystemId)
			}

			name, ok := resultMap["name"].(string)
			if !ok {
				log.Error(errors.New(fmt.Sprintf("starSystem %d %v cast to string failed", starSystem.SystemId, resultMap["name"])))
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
			log.Errorf(err, "starSystem %d save to DB failed", starSystem.SystemId)
		}
	}
}
