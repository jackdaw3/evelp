package sde

import (
	"evelp/model"
	"io/ioutil"
	"sort"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type FactionsInit struct {
	FilePath string
	factions *model.Factions
}

func (f *FactionsInit) Refresh() error {
	log.Info("Start load factions", f.FilePath)
	if err := f.load(); err != nil {
		return err
	}
	log.Info("Load ", f.FilePath, " finished.")

	log.Info("Start save factions to DB.")
	if err := model.SaveFactions(f.factions); err != nil {
		return err
	}
	log.Info("Factions have saved to DB.")

	return nil
}

func (f *FactionsInit) load() error {
	file, err := ioutil.ReadFile(f.FilePath)
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
