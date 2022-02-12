package service

import (
	"evelp/log"
	"evelp/model"
	"evelp/util/cache"
	"strconv"
)

const history = "history"

type ItemHistoryService struct {
	itemId   int
	regionId int
}

func NewItemHistoryService(itemId int, regionId int) *ItemHistoryService {
	return &ItemHistoryService{itemId, regionId}
}

func (h *ItemHistoryService) History() (*model.ItemHistorys, error) {
	var itemHistorys model.ItemHistorys

	key := cache.Key(history, strconv.Itoa(h.regionId), strconv.Itoa(h.itemId))
	exist, err := cache.Exist(key)
	if err != nil {
		return nil, err
	}

	if exist {
		if err := cache.Get(key, &itemHistorys); err != nil {
			return nil, err
		}
	} else {
		log.Warnf("key %v not exist in redis", key)
		return nil, nil
	}

	return &itemHistorys, nil
}
