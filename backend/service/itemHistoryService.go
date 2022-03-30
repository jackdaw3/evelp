package service

import (
	"evelp/dto"
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

func (ih *ItemHistoryService) AverageVolume(days int) (int64, error) {
	itemHistorys, err := ih.historyFromCache()
	if err != nil {
		return 0, err
	}

	volume := itemHistorys.AverageVolume(days)
	return volume, nil
}

func (ih *ItemHistoryService) History() (*dto.ItemHistoryDTOs, error) {
	historys, err := ih.historyFromCache()
	if err != nil {
		return nil, err
	}

	var historyDTOs dto.ItemHistoryDTOs
	for _, history := range *historys {
		var histroyDTO dto.ItemHistoryDTO
		histroyDTO.ItemId = history.ItemId
		histroyDTO.Average = history.Average
		histroyDTO.Highest = history.Highest
		histroyDTO.Lowest = history.Lowest
		histroyDTO.OrderCount = history.OrderCount
		histroyDTO.Volume = history.Volume
		histroyDTO.Date = history.Date
		historyDTOs = append(historyDTOs, &histroyDTO)
	}

	historyDTOs.GenerateHistory()
	return &historyDTOs, nil
}

func (ih *ItemHistoryService) historyFromCache() (*model.ItemHistorys, error) {
	var itemHistorys model.ItemHistorys

	var key string
	if ih.isBluePrint {
		bluePrint, err := model.GetBluePrint(ih.itemId)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to get blue print %d", ih.itemId)
		}
		if len(bluePrint.Products) == 0 {
			return nil, errors.Errorf("offer %d's bluePrint %d have no product", ih.itemId, bluePrint.BlueprintId)
		}
		key = cache.Key(history, strconv.Itoa(ih.regionId), strconv.Itoa(bluePrint.Products[0].ItemId))
	} else {
		key = cache.Key(history, strconv.Itoa(ih.regionId), strconv.Itoa(ih.itemId))

	}

	if err := cache.Get(key, &itemHistorys); err != nil {
		return nil, err
	}

	return &itemHistorys, nil
}
