package service

import (
	"encoding/json"
	"evelp/model"
	"evelp/util/cache"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestAverageVolume(t *testing.T) {
	defer monkey.UnpatchAll()
	itemHistorys := mockHistoryData()
	mockHistory(itemHistorys)

	histroyService := NewItemHistoryService(34, 10000002, false)
	volume, err := histroyService.AverageVolume(2)
	assert.NoError(t, err)
	assert.Equal(t, (6358751950+3910212408)/2, int(volume))

	volume, err = histroyService.AverageVolume(7)
	assert.NoError(t, err)
	assert.Equal(t, (6358751950+3910212408)/2, int(volume))
}

func TestHistory(t *testing.T) {
	defer monkey.UnpatchAll()
	itemHistorys := mockHistoryData()
	mockHistory(itemHistorys)

	histroyService := NewItemHistoryService(34, 10000002, false)
	history, err := histroyService.History()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(*history))
}

func mockHistoryData() *model.ItemHistorys {
	return &model.ItemHistorys{
		&model.ItemHistory{ItemId: 34,
			Average:    6.81,
			Highest:    6.96,
			Lowest:     6.7,
			OrderCount: 2166,
			Volume:     3910212408,
			Date:       "2021-01-01",
		},
		&model.ItemHistory{ItemId: 34,
			Average:    6.54,
			Highest:    6.93,
			Lowest:     6.45,
			OrderCount: 2513,
			Volume:     6358751950,
			Date:       "2021-01-02",
		},
	}
}

func mockHistory(itemHistorys *model.ItemHistorys) {
	monkey.Patch(cache.Get, func(key string, dest interface{}) error {
		val, err := json.Marshal(itemHistorys)
		if err != nil {
			return err
		}
		json.Unmarshal([]byte(val), dest)
		return nil
	})
}
