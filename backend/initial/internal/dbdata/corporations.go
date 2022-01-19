package dbdata

import (
	"evelp/model"
	"io/ioutil"
	"sort"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type corporationsData struct {
	filePath     string
	corporations *model.Corporations
}

func (c *corporationsData) Refresh() error {
	log.Info("start load corporations", c.filePath)
	if err := c.load(); err != nil {
		return err
	}
	log.Info("load ", c.filePath, " finished")

	log.Info("start save corporations to DB")
	if err := model.SaveCorporations(c.corporations); err != nil {
		return err
	}
	log.Info("corporations have saved to DB")

	return nil
}

func (c *corporationsData) load() error {
	file, err := ioutil.ReadFile(c.filePath)
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
	c.corporations = &corporations

	return nil
}
