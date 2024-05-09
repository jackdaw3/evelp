package service

import (
	"evelp/dto"
	"evelp/model"
)

type ItemService struct {
	lang string
}

func NewItemService(lang string) *ItemService {
	return &ItemService{lang}
}

func (is *ItemService) Item(itemId int) (*dto.ItemDTO, error) {
	var itemDTO dto.ItemDTO

	item, err := model.GetItem(itemId)
	if err != nil {
		return nil, err
	}

	itemDTO.ItemId = item.ItemId
	itemDTO.ItemName = item.Name.Lang(is.lang)

	return &itemDTO, nil
}

func (is *ItemService) ItemDetail(itemId int) (*dto.ItemDetailDTO, error) {
	var itemDetailDTO dto.ItemDetailDTO

	item, err := model.GetItem(itemId)
	if err != nil {
		return nil, err
	}

	itemDetailDTO.ItemId = item.ItemId
	itemDetailDTO.ItemName = item.Name.Lang(is.lang)
	itemDetailDTO.Description = item.Description.Lang(is.lang)
	itemDetailDTO.Volume = item.Volume

	return &itemDetailDTO, nil
}
