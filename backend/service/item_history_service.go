package service

import (
	"evelp/model"
	"evelp/util/cache"
	"strconv"

	"github.com/pkg/errors"
)

type ItemHistoryService struct {
	itemId   int
	regionId int
}

func NewItemHistoryService(itemId int, regionId int) *ItemHistoryService {
	return &ItemHistoryService{itemId, regionId}
}

func (h *ItemHistoryService) History() (*model.ItemHistorys, error) {
	var itemHistorys model.ItemHistorys
	key := cache.Key(order, strconv.Itoa(h.regionId), strconv.Itoa(h.itemId))
	if err := cache.Get(key, &itemHistorys); err != nil {
		return nil, errors.WithMessagef(err, "get itemHistorys %s cache error", key)
	}
	return &itemHistorys, nil
}
