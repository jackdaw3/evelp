package service

import (
	"evelp/model"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestFactions(t *testing.T) {
	defer monkey.UnpatchAll()
	factions, coporations := mockFactionsData()
	mockFactions(factions, coporations)

	factionService := NewFactionService("en")
	factionDTOs, err := factionService.Factions()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(*factions))
	assert.Equal(t, 2, len((*factionDTOs)[0].Corporations))
	assert.Equal(t, (*coporations)[0].Name.En, (*factionDTOs)[0].Corporations[0].CorporationName)
}

func mockFactions(factions *model.Factions, coporations *model.Corporations) {
	monkey.Patch(model.GetFactions, func() (*model.Factions, error) {
		return factions, nil
	})

	monkey.Patch(model.GetCorporationsByFaction, func(int) (*model.Corporations, error) {
		return coporations, nil
	})
}

func mockFactionsData() (*model.Factions, *model.Corporations) {
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
	return factions, coporations
}
