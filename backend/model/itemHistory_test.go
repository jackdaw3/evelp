package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAverageVolume(t *testing.T) {
	itemHistorys := &ItemHistorys{
		&ItemHistory{Date: "2021-01-01", Volume: 3910212408},
		&ItemHistory{Date: "2021-01-02", Volume: 6358751950},
		&ItemHistory{Date: "2021-01-03", Volume: 7426253889},
		&ItemHistory{Date: "2021-01-04", Volume: 4974822193},
		&ItemHistory{Date: "2021-01-05", Volume: 6614498676},
		&ItemHistory{Date: "2021-01-06", Volume: 4342137765},
		&ItemHistory{Date: "2021-01-07", Volume: 5576082901},
		&ItemHistory{Date: "2021-01-08", Volume: 3995459669},
		&ItemHistory{Date: "2021-01-09", Volume: 7086188858},
		&ItemHistory{Date: "2021-01-10", Volume: 6539802438},
	}

	assert.Equal(t, int64(5873816988), itemHistorys.AverageVolume(3))
	assert.Equal(t, int64(5589856071), itemHistorys.AverageVolume(7))
	assert.Equal(t, int64(5682421074), itemHistorys.AverageVolume(15))

}
