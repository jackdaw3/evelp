package service

import (
	"evelp/model"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

var (
	factions = &model.Factions{
		model.Faction{FactionId: 500001,
			Name: model.Name{En: "Caldari State"},
		},
	}
	coporations = &model.Corporations{
		model.Corporation{CorporationId: 1000002,
			FactionId: 500001,
			Name:      model.Name{En: "CBD Corporation"},
		},
		model.Corporation{CorporationId: 1000003,
			FactionId: 500001,
			Name:      model.Name{En: "Prompt Delivery"},
		},
	}
)

func TestFactions(t *testing.T) {
	defer monkey.UnpatchAll()
	en := "en"
	factionService := NewFactionService(en)

	monkey.Patch(model.GetFactions, func() (*model.Factions, error) {
		return factions, nil
	})

	monkey.Patch(model.GetCorporationsByFaction, func(int) (*model.Corporations, error) {
		return coporations, nil
	})

	factions, err := factionService.Factions()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(*factions))
	assert.Equal(t, 2, len((*factions)[0].Corporations))
	assert.Equal(t, (*coporations)[0].Name.En, (*factions)[0].Corporations[0].CorporationName)
}
