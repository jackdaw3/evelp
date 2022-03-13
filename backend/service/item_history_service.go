package service

import (
	"evelp/model"
	"evelp/util/cache"
	"strconv"

	"github.com/pkg/errors"
)

const history = "history"

type ItemHistoryService struct {
	itemId      int
	regionId    int
	isBluePrint bool
}

func NewItemHistoryService(itemId int, regionId int, isBluePrint bool) *ItemHistoryService {
	return &ItemHistoryService{itemId, regionId, isBluePrint}
}

func (h *ItemHistoryService) History() (*model.ItemHistorys, error) {
	var itemHistorys model.ItemHistorys

	var key string
	if h.isBluePrint {
		bluePrint := model.GetBluePrint(h.itemId)
		if len(bluePrint.Products) == 0 {
			return nil, errors.Errorf("offer %d's bluePrint %d have no product", h.itemId, bluePrint.BlueprintId)
		}
		key = cache.Key(history, strconv.Itoa(h.regionId), strconv.Itoa(bluePrint.Products[0].ItemId))
	} else {
		key = cache.Key(history, strconv.Itoa(h.regionId), strconv.Itoa(h.itemId))

	}

	if err := cache.Get(key, &itemHistorys); err != nil {
		return nil, err
	}

	return &itemHistorys, nil
}
