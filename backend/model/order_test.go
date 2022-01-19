package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	orders             Orders
	multipleItemOrders Orders
	buyOrders          Orders
	sellOrders         Orders
)

var (
	scope1 float64 = 0.01
	scope2 float64 = 0.05
	scope3 float64 = 0.1
	scope4 float64 = 0.2
	scope5 float64 = 0.5
)

func setUp() {
	time1, _ := time.Parse(time.RFC3339, "2022-01-13T22:51:59Z")
	time2, _ := time.Parse(time.RFC3339, "2022-01-11T18:36:34Z")
	time3, _ := time.Parse(time.RFC3339, "2022-01-08T13:38:28Z")
	time4, _ := time.Parse(time.RFC3339, "2021-12-15T11:52:04Z")
	time5, _ := time.Parse(time.RFC3339, "2021-12-15T11:54:30Z")
	time6, _ := time.Parse(time.RFC3339, "2022-01-14T04:58:10Z")

	order1 := Order{6173392220, 28758, time1, 90, 30000142, 26350000, 28, 30, true, time.Now()}
	order2 := Order{6171724721, 28758, time2, 90, 30000142, 27850000, 3, 30, true, time.Now()}
	order3 := Order{6169089210, 28758, time3, 90, 30000142, 27760000, 5, 10, true, time.Now()}

	order4 := Order{6150131000, 28758, time4, 90, 30000142, 38500000, 92, 100, false, time.Now()}
	order5 := Order{6150132220, 28758, time5, 90, 30000142, 38100000, 9, 15, false, time.Now()}
	order6 := Order{6173556403, 28758, time6, 90, 30000142, 36730000, 12, 50, false, time.Now()}
	anotherItemOrder := Order{6150131000, 28759, time4, 90, 30000142, 38500000, 92, 100, false, time.Now()}

	orders = Orders{order5, order1, order4, order2, order3, order6}
	buyOrders = Orders{order1, order2, order3}
	sellOrders = Orders{order4, order5, order6}
	multipleItemOrders = Orders{order1, order4, anotherItemOrder}
}
func TestGetHighestBuyPrice(t *testing.T) {
	setUp()

	invaidHighestPrice1, err := multipleItemOrders.GetHighestBuyPrice(scope1)
	assert.Zero(t, invaidHighestPrice1)
	assert.Equal(t, "orders have multiple itemIds", err.Error())

	invaidHighestPrice2, err := sellOrders.GetHighestBuyPrice(scope1)
	assert.Zero(t, invaidHighestPrice2)
	assert.NoError(t, err)

	highestPrice1, err := orders.GetHighestBuyPrice(scope1)
	assert.Equal(t, float64(27850000), highestPrice1)
	assert.NoError(t, err)

	highestPrice2, err := orders.GetHighestBuyPrice(scope2)
	assert.Equal(t, float64(27850000), highestPrice2)
	assert.NoError(t, err)

	highestPrice3, err := orders.GetHighestBuyPrice(scope3)
	assert.Equal(t, float64(27827500), highestPrice3)
	assert.NoError(t, err)

	highestPrice4, err := orders.GetHighestBuyPrice(scope4)
	assert.Equal(t, float64(27793750), highestPrice4)
	assert.NoError(t, err)

	highestPrice5, err := orders.GetHighestBuyPrice(scope5)
	assert.Equal(t, int64(26991666), int64(highestPrice5))
	assert.NoError(t, err)
}

func TestGetLowestSellPrice(t *testing.T) {
	setUp()

	invaidLowestPrice1, err := multipleItemOrders.GetLowestSellPrice(scope1)
	assert.Zero(t, invaidLowestPrice1)
	assert.Equal(t, "orders have multiple itemIds", err.Error())

	invaidLowestPrice2, err := buyOrders.GetLowestSellPrice(scope1)
	assert.Zero(t, invaidLowestPrice2)
	assert.NoError(t, err)

	lowestPrice1, err := orders.GetLowestSellPrice(scope1)
	assert.Equal(t, float64(36730000), lowestPrice1)
	assert.NoError(t, err)

	lowestPrice2, err := orders.GetLowestSellPrice(scope2)
	assert.Equal(t, float64(36730000), lowestPrice2)
	assert.NoError(t, err)

	lowestPrice3, err := orders.GetLowestSellPrice(scope3)
	assert.Equal(t, float64(36730000), lowestPrice3)
	assert.NoError(t, err)

	lowestPrice4, err := orders.GetLowestSellPrice(scope4)
	assert.Equal(t, float64(37420000), lowestPrice4)
	assert.NoError(t, err)

	lowestPrice5, err := orders.GetLowestSellPrice(scope5)
	assert.Equal(t, int64(38064210), int64(lowestPrice5))
	assert.NoError(t, err)
}
