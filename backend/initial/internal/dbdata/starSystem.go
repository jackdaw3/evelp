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

type starSystemData struct {
	starSystems *model.StarSystems
}

func (s *starSystemData) Refresh() error {
	log.Infof("start to load starSystems from %s", global.Conf.Data.Remote.Address)
	s.getAllStarSystems()
	sort.Sort(s.starSystems)

	for _, starSystem := range *s.starSystems {
		exist, err := starSystem.IsExist()
		if err != nil {
			log.Errorf(err, "failed to check starSystem %d exist", starSystem.SystemId)
		}

		if exist {
			valid, err := starSystem.IsVaild()
			if err != nil {
				log.Errorf(err, "failed to check starSystem %d valid", starSystem.SystemId)
			}

			if valid {
				continue
			}
		}

		wg.Add(1)
		global.Ants.Submit(s.getStarSystem(starSystem, &wg))
	}

	wg.Wait()
	log.Info("starSystems loaded and saved to DB")

	return nil
}

func (s *starSystemData) getAllStarSystems() {
	req := fmt.Sprintf("%s/universe/systems/?datasource=%s",
		global.Conf.Data.Remote.Address,
		global.Conf.Data.Remote.DataSource,
	)

	resp, err := net.GetWithRetries(client, req)
	if err != nil {
		log.Errorf(err, "failed to get starSystems")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf(err, "failed to get starSystems' body")
	}

	var idArray []int

	if err = json.Unmarshal(body, &idArray); err != nil {
		log.Errorf(err, "failed to unmarshal starSystems json")
	}

	var starSystems model.StarSystems
	for _, id := range idArray {
		var starSystem model.StarSystem
		starSystem.SystemId = id
		starSystems = append(starSystems, &starSystem)
	}
	s.starSystems = &starSystems
}

func (s *starSystemData) getStarSystem(starSystem *model.StarSystem, wg *sync.WaitGroup) func() {
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
				log.Errorf(err, "failed to get starSystem %d", starSystem.SystemId)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf(err, "failed to get starSystem %d's body", starSystem.SystemId)
			}

			var resultMap map[string]interface{}

			if err = json.Unmarshal(body, &resultMap); err != nil {
				log.Errorf(err, "failed to unmarshal starSystem %d json", starSystem.SystemId)
			}

			name, ok := resultMap["name"].(string)
			if !ok {
				log.Error(errors.New(fmt.Sprintf("failed to cast starSystem %d %v to string ", starSystem.SystemId, resultMap["name"])))
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
			log.Errorf(err, "failed to starSystem %d save to DB", starSystem.SystemId)
		}
	}
}
