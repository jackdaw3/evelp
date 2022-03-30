package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHistory(t *testing.T) {
	list := &ItemHistoryDTOs{
		&ItemHistoryDTO{Average: 50, Highest: 60, Lowest: 40},
		&ItemHistoryDTO{Average: 60, Highest: 70, Lowest: 50},
		&ItemHistoryDTO{Average: 70, Highest: 80, Lowest: 60},
		&ItemHistoryDTO{Average: 20, Highest: 30, Lowest: 10},
		&ItemHistoryDTO{Average: 30, Highest: 40, Lowest: 20},
	}

	list.GenerateHistory()

	assert.Equal(t, (*list)[0].Average5d, float64(50))
	assert.Equal(t, (*list)[1].Average5d, float64(55))
	assert.Equal(t, (*list)[2].Average5d, float64(60))
	assert.Equal(t, (*list)[3].Average5d, float64(50))
	assert.Equal(t, (*list)[4].Average5d, float64(46))

	assert.Equal(t, (*list)[0].Average20d, float64(50))
	assert.Equal(t, (*list)[1].Average20d, float64(55))
	assert.Equal(t, (*list)[2].Average20d, float64(60))
	assert.Equal(t, (*list)[3].Average20d, float64(50))
	assert.Equal(t, (*list)[4].Average20d, float64(46))

	assert.Equal(t, (*list)[0].Highest5d, float64(60))
	assert.Equal(t, (*list)[1].Highest5d, float64(70))
	assert.Equal(t, (*list)[2].Highest5d, float64(80))
	assert.Equal(t, (*list)[3].Highest5d, float64(80))
	assert.Equal(t, (*list)[4].Highest5d, float64(80))

	assert.Equal(t, (*list)[0].Lowest5d, float64(40))
	assert.Equal(t, (*list)[1].Lowest5d, float64(40))
	assert.Equal(t, (*list)[2].Lowest5d, float64(40))
	assert.Equal(t, (*list)[3].Lowest5d, float64(10))
	assert.Equal(t, (*list)[4].Lowest5d, float64(10))
}
