package dbdata

import (
	"evelp/log"
	"evelp/model"
	"io/ioutil"
	"sort"

	"gopkg.in/yaml.v2"
)

type corporationData struct {
	filePath     string
	corporations *model.Corporations
}

func (cd *corporationData) Refresh() error {
	log.Info("start to load corporations", cd.filePath)
	if err := cd.loadCorporations(); err != nil {
		return err
	}
	log.Infof("loading %s completed", cd.filePath)

	log.Info("start to save corporations to DB")
	if err := model.SaveCorporations(cd.corporations); err != nil {
		return err
	}
	log.Info("corporations saved to DB")

	return nil
}

func (cd *corporationData) loadCorporations() error {
	file, err := ioutil.ReadFile(cd.filePath)
	if err != nil {
		return err
	}

	data := make(map[int]model.Corporation)
	if err := yaml.Unmarshal(file, &data); err != nil {
		return err
	}

	var corporations model.Corporations
	for id, corporation := range data {
		corporation.CorporationId = id
		corporations = append(corporations, corporation)
	}
	sort.Sort(corporations)
	cd.corporations = &corporations

	return nil
}
