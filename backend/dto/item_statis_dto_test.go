package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUnitProfit(t *testing.T) {
	orderDTOWrappers := OrderDTOWrappers{
		&OrderDTOWrapper{
			OrderDTO: OrderDTO{OrderId: 6206480166,
				ItemId:       28813,
				VolumeRemain: 1,
				VolumeTotal:  1,
				IsBuyOrder:   true,
				Price:        755700000,
				LastUpdated:  time.Now().GoString(),
			},
			Income: 755700000,
			Cost:   242393333,
			Profit: 513306666,
		},
		&OrderDTOWrapper{
			OrderDTO: OrderDTO{OrderId: 6205771440,
				ItemId:       28813,
				VolumeRemain: 2,
				VolumeTotal:  2,
				IsBuyOrder:   true,
				Price:        755500000,
				LastUpdated:  time.Now().GoString(),
			},
			Income: 1511000000,
			Cost:   484786666,
			Profit: 1026213333,
		},
	}

	itemStatis := new(ItemStatisDTO)
	itemStatis.Orderwrappers = orderDTOWrappers
	itemStatis.GenerateUnitProfit(321300)
	assert.Equal(t, 1597, int(itemStatis.AveUnitProfit))
}
