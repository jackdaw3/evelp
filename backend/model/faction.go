package model

import (
	"evelp/config/global"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

type Faction struct {
	FactionId int  `gorm:"type:int;not null;primary_key;autoIncrement:false"`
	Name      Name `gorm:"type:text" yaml:"nameID"`
}

type Factions []Faction

func (factions Factions) Len() int { return len(factions) }

func (factions Factions) Less(i, j int) bool { return factions[i].FactionId < factions[j].FactionId }

func (factions Factions) Swap(i, j int) { factions[i], factions[j] = factions[j], factions[i] }

func GetFaction(id int) (*Faction, error) {
	var faction Faction
	result := global.DB.First(&faction, id)
	return &faction, result.Error
}

func GetFactions() (*Factions, error) {
	var factions Factions
	result := global.DB.Find(&factions)
	return &factions, result.Error
}

func SaveFaction(faction *Faction) error {
	if err := global.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&faction).Error; err != nil {
		return errors.Wrapf(err, "failed to save faction %d to DB", faction.FactionId)
	}
	return nil
}

func SaveFactions(factions *Factions) error {
	for _, faction := range *factions {
		if err := SaveFaction(&faction); err != nil {
			return err
		}
	}
	return nil
}
