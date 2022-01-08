package sde

import (
	"evelp/model"
	"io/ioutil"
	"sort"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type TypeIDsInit struct {
	FilePath string
	items    *model.Items
}

func (t *TypeIDsInit) Refresh() error {
	log.Info("Start load items", t.FilePath)
	if err := t.load(); err != nil {
		return err
	}
	log.Info("Load ", t.FilePath, " finished.")

	log.Info("Save types to DB.")
	if err := model.SaveItems(t.items); err != nil {
		return err
	}
	log.Info("Types have saved to DB.")

	return nil
}

func (t *TypeIDsInit) load() error {
	file, err := ioutil.ReadFile(t.FilePath)
	if err != nil {
		return err
	}

	data := make(map[int]model.Item)
	if err := yaml.Unmarshal(file, &data); err != nil {
		return err
	}

	var items model.Items
	for id, item := range data {
		item.ItemId = id
		items = append(items, item)
	}
	sort.Sort(items)
	t.items = &items

	return nil
}
