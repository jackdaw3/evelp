package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCost(t *testing.T) {
	materailDTOs := MatertialDTOs{MaterialDTO{Cost: 1000}, MaterialDTO{Cost: 2000}, MaterialDTO{Cost: 3000}}

	assert.Equal(t, float64(6000), materailDTOs.Cost())
}
