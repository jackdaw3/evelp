package service

import (
	"evelp/dto"
	"evelp/log"
	"evelp/model"
	"sort"
)

type FactionSerivce struct {
	lang string
}

func NewFactionService(lang string) *FactionSerivce {
	return &FactionSerivce{lang}
}

func (f *FactionSerivce) Factions() (*dto.FactionDTOs, error) {
	var factionDTOs dto.FactionDTOs

	factions, err := model.GetFactions()
	if err != nil {
		return nil, err
	}

	for _, faction := range *factions {
		var factionDTO dto.FactionDTO
		factionDTO.FactionId = faction.FactionId
		factionDTO.FactionName = faction.Name.Lang(f.lang)

		var corporationDTOs dto.CorporationDTOs
		corportations, err := model.GetCorporationsByFaction(factionDTO.FactionId)
		if err != nil {
			log.Errorf(err, "failed to get faction %d's corporation list", factionDTO.FactionId)
			continue
		}

		for _, corportation := range *corportations {
			var corporationDTO dto.CorporationDTO
			corporationDTO.CorporationId = corportation.CorporationId
			corporationDTO.CorporationName = corportation.Name.Lang(f.lang)
			corporationDTOs = append(corporationDTOs, &corporationDTO)
		}
		sort.Sort(corporationDTOs)
		factionDTO.Corporations = corporationDTOs

		factionDTOs = append(factionDTOs, &factionDTO)
	}

	sort.Sort(factionDTOs)
	return &factionDTOs, nil
}
