package service

import (
	"evelp/model"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestCorporation(t *testing.T) {
	defer monkey.UnpatchAll()
	corporation := mockCorporationData()
	mockCorporation(corporation)

	corporationService := NewCorporationSerivce(corporation.CorporationId, "en")
	corporationDTO, err := corporationService.Corporation()
	assert.NoError(t, err)
	assert.Equal(t, corporation.Name.En, corporationDTO.CorporationName)
	assert.Equal(t, corporation.CorporationId, corporationDTO.CorporationId)

}

func mockCorporation(corporation *model.Corporation) {
	monkey.Patch(model.GetCorporation, func(int) (*model.Corporation, error) {
		return corporation, nil
	})
}

func mockCorporationData() *model.Corporation {
	return &model.Corporation{CorporationId: 1000002,
		FactionId: 500001,
		Name:      model.Name{En: "CBD Corporation"},
	}
}
