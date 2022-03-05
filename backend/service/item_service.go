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
	itemDTO.ItemName = item.Name.Val(is.lang)

	
	return &itemDTO, nil
}
