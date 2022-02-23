package service

import (
	"encoding/json"
	"evelp/model"
	"evelp/util/cache"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

var (
	itemHistorys = &model.ItemHistorys{
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
)

func TestHistory(t *testing.T) {
	defer monkey.UnpatchAll()
	histroyService := NewItemHistoryService(34, 10000002)

	monkey.Patch(cache.Get, func(key string, dest interface{}) error {
		val, err := json.Marshal(itemHistorys)
		if err != nil {
			return err
		}
		json.Unmarshal([]byte(val), dest)
		return nil
	})

	historys, err := histroyService.History()
	assert.NoError(t, err)
	assert.Equal(t, len(*itemHistorys), len(*historys))

}
