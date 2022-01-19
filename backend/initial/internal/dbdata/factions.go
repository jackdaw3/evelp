package dbdata

import (
	"evelp/model"
	"io/ioutil"
	"sort"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type factionsData struct {
	filePath string
	factions *model.Factions
}

func (f *factionsData) Refresh() error {
	log.Info("start load factions", f.filePath)
	if err := f.load(); err != nil {
		return err
	}
	log.Info("load ", f.filePath, " finished")

	log.Info("start save factions to DB")
	if err := model.SaveFactions(f.factions); err != nil {
		return err
	}
	log.Info("factions have saved to DB")

	return nil
}

func (f *factionsData) load() error {
	file, err := ioutil.ReadFile(f.filePath)
	if err != nil {
		return err
	}

	data := make(map[int]model.Faction)
	if err := yaml.Unmarshal(file, &data); err != nil {
		return err
	}

	var factions model.Factions
	for id, faction := range data {
		faction.FactionId = id
		factions = append(factions, faction)
	}
	sort.Sort(factions)
	f.factions = &factions

	return nil
}
