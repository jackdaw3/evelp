package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var orders Orders

func setUp() {
	time1, _ := time.Parse(time.RFC3339, "2022-01-13T22:51:59Z")
	time2, _ := time.Parse(time.RFC3339, "2022-01-11T18:36:34Z")
	time3, _ := time.Parse(time.RFC3339, "2022-01-08T13:38:28Z")
	time4, _ := time.Parse(time.RFC3339, "2021-12-15T11:52:04Z")
	time5, _ := time.Parse(time.RFC3339, "2021-12-15T11:54:30Z")
	time6, _ := time.Parse(time.RFC3339, "2022-01-14T04:58:10Z")

	orders = Orders{
		Order{6150132220, 28758, time5, 90, 30000142, 38100000, 9, 15, false},
		Order{6173392220, 28758, time1, 90, 30000142, 26350000, 28, 30, true},
		Order{6150131000, 28758, time4, 90, 30000142, 38500000, 92, 100, false},
		Order{6171724721, 28758, time2, 90, 30000142, 27850000, 3, 30, true},
		Order{6169089210, 28758, time3, 90, 30000142, 27760000, 5, 10, true},
		Order{6173556403, 28758, time6, 90, 30000142, 36730000, 12, 50, false},
	}
}
func TestGetHighestBuyPrice(t *testing.T) {
	setUp()

	var scope1 float64 = 0.01
	highestPrice1, _ := orders.GetHighestBuyPrice(scope1)
	assert.Equal(t, float64(27850000), highestPrice1)

	var scope2 float64 = 0.05
	highestPrice2, _ := orders.GetHighestBuyPrice(scope2)
	assert.Equal(t, float64(27850000), highestPrice2)

	var scope3 float64 = 0.1
	highestPrice3, _ := orders.GetHighestBuyPrice(scope3)
	assert.Equal(t, float64(27827500), highestPrice3)

	var scope4 float64 = 0.2
	highestPrice4, _ := orders.GetHighestBuyPrice(scope4)
	assert.Equal(t, float64(27793750), highestPrice4)

	var scope5 float64 = 0.5
	highestPrice5, _ := orders.GetHighestBuyPrice(scope5)
	assert.Equal(t, int64(26991666), int64(highestPrice5))
}

func TestGetLowestSellPrice(t *testing.T) {
	setUp()

	var scope1 float64 = 0.01
	lowestPrice1, _ := orders.GetLowestSellPrice(scope1)
	assert.Equal(t, float64(36730000), lowestPrice1)

	var scope2 float64 = 0.05
	lowestPrice2, _ := orders.GetLowestSellPrice(scope2)
	assert.Equal(t, float64(36730000), lowestPrice2)

	var scope3 float64 = 0.1
	lowestPrice3, _ := orders.GetLowestSellPrice(scope3)
	assert.Equal(t, float64(36730000), lowestPrice3)

	var scope4 float64 = 0.2
	lowestPrice4, _ := orders.GetLowestSellPrice(scope4)
	assert.Equal(t, float64(37420000), lowestPrice4)

	var scope5 float64 = 0.5
	lowestPrice5, _ := orders.GetLowestSellPrice(scope5)
	assert.Equal(t, int64(38064210), int64(lowestPrice5))
}
