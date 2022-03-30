package dbdata

import (
	"evelp/log"
	"evelp/model"
	"io/ioutil"
	"sort"

	"gopkg.in/yaml.v2"
)

type typeIDData struct {
	filePath string
	items    *model.Items
}

func (t *typeIDData) Refresh() error {
	log.Info("start to load items", t.filePath)
	if err := t.load(); err != nil {
		return err
	}
	log.Infof("loading %s completed", t.filePath)

	log.Info("start to save types to DB.")
	if err := model.SaveItems(t.items); err != nil {
		return err
	}
	log.Info("types saved to DB")

	return nil
}

func (t *typeIDData) load() error {
	file, err := ioutil.ReadFile(t.filePath)
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
