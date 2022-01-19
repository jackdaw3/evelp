package dbdata

import (
	"evelp/model"
	"io/ioutil"
	"sort"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type typeIdsData struct {
	filePath string
	items    *model.Items
}

func (t *typeIdsData) Refresh() error {
	log.Info("Start load items", t.filePath)
	if err := t.load(); err != nil {
		return err
	}
	log.Info("Load ", t.filePath, " finished.")

	log.Info("Start save types to DB.")
	if err := model.SaveItems(t.items); err != nil {
		return err
	}
	log.Info("Types have saved to DB.")

	return nil
}

func (t *typeIdsData) load() error {
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
