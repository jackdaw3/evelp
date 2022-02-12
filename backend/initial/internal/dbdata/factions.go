package dbdata

import (
	"evelp/log"
	"evelp/model"
	"io/ioutil"
	"sort"

	"gopkg.in/yaml.v2"
)

type factionsData struct {
	filePath string
	factions *model.Factions
}

func (fd *factionsData) Refresh() error {
	log.Info("start load factions", fd.filePath)
	if err := fd.load(); err != nil {
		return err
	}
	log.Info("load ", fd.filePath, " finished")

	log.Info("start save factions to DB")
	if err := model.SaveFactions(fd.factions); err != nil {
		return err
	}
	log.Info("factions have saved to DB")

	return nil
}

func (fd *factionsData) load() error {
	file, err := ioutil.ReadFile(fd.filePath)
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
	fd.factions = &factions

	return nil
}
