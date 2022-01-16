package sde

import (
	"evelp/model"
	"io/ioutil"
	"sort"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type CorporationsInit struct {
	FilePath     string
	corporations *model.Corporations
}

func (c *CorporationsInit) Refresh() error {
	log.Info("Start load corporations", c.FilePath)
	if err := c.load(); err != nil {
		return err
	}
	log.Info("Load ", c.FilePath, " finished.")

	log.Info("Start save corporations to DB.")
	if err := model.SaveCorporations(c.corporations); err != nil {
		return err
	}
	log.Info("Corporations have saved to DB.")

	return nil
}

func (c *CorporationsInit) load() error {
	file, err := ioutil.ReadFile(c.FilePath)
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
