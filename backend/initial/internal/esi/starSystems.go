package esi

import (
	"encoding/json"
	"evelp/config/global"
	"evelp/model"
	"evelp/util/netUtil"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
)

type StarSystemsInit struct {
	starSystems *model.StarSystems
}

func (s *StarSystemsInit) Refresh() error {
	log.Infof("Start load starSystems from %s.", global.Conf.Data.RemoteDataAddress)
	s.getAllStarSystems()
	sort.Sort(s.starSystems)

	for _, starSystem := range *s.starSystems {
		exist, err := model.IsStarSystemExist(starSystem.SystemId)
		if err != nil {
			log.Errorf("Check starSystem %d exist failed %s.", starSystem.SystemId, err)
		}

		if exist {
			continue
		}

		wg.Add(1)
		acquireSem(weigth)
		go getStarSystem(starSystem)
	}
	wg.Wait()
	log.Info("StarSystems loaded and have saved to DB..")

	return nil
}

func (s *StarSystemsInit) getAllStarSystems() {
	req := fmt.Sprintf("%s/universe/systems/?datasource=%s", global.Conf.Data.RemoteDataAddress, global.Conf.Data.RemoteDataSource)

	body, err := netUtil.GetWithRetries(client, req)
	if err != nil {
		log.Errorf("Get starSystems failed: %s", err.Error())
	}

	var idArray []int

	if err = json.Unmarshal(body, &idArray); err != nil {
		log.Errorf("Unmarshal starSystems json failed: %s", err.Error())
	}

	var starSystems model.StarSystems
	for _, id := range idArray {
		var starSystem model.StarSystem
		starSystem.SystemId = id
		starSystems = append(starSystems, &starSystem)
	}
	s.starSystems = &starSystems
}

func getStarSystem(starSystem *model.StarSystem) {
	defer wg.Done()
	defer sem.Release(weigth)

	for _, lang := range langs {
		req := fmt.Sprintf("%s/universe/systems/%d/?datasource=%s&language=%s", global.Conf.Data.RemoteDataAddress, starSystem.SystemId, global.Conf.Data.RemoteDataSource, lang)

		body, err := netUtil.GetWithRetries(client, req)
		if err != nil {
			log.Errorf("Get starSystem %d failed: %s", starSystem.SystemId, err.Error())
		}

		var resultMap map[string]interface{}

		if err = json.Unmarshal(body, &resultMap); err != nil {
			log.Errorf("Unmarshal starSystem %d json failed: %s", starSystem.SystemId, err.Error())
		}

		name, ok := resultMap["name"].(string)
		if !ok {
			log.Errorf("StarSystem %d %v cast to string failed.", starSystem.SystemId, resultMap["name"])
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
		log.Errorf("StarSystem %d failed to save to DB.", starSystem.SystemId)
	}
}
