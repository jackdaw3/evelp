package service

import (
	"evelp/dto"
	"evelp/model"
)

type CorporationSerivce struct {
	corporationId int
	lang          string
}

func NewCorporationSerivce(corporationId int, lang string) *CorporationSerivce {
	return &CorporationSerivce{corporationId, lang}
}

func (c *CorporationSerivce) Corporation() (*dto.CorporationDTO, error) {
	corporation, err := model.GetCorporation(c.corporationId)
	if err != nil {
		return nil, err
	}

	var corporationDTO dto.CorporationDTO
	corporationDTO.CorporationId = c.corporationId
	corporationDTO.CorporationName = corporation.Name.Lang(c.lang)

	return &corporationDTO, nil
}
