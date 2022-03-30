package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSalaIndex(t *testing.T) {
	offerDTO := new(OfferDTO)
	offerDTO.Price = k1 * 2
	offerDTO.Volume = 128
	offerDTO.UnitProfit = 2500

	offerDTO.GenerateSaleIndex()

	assert.Equal(t, 226, offerDTO.SaleIndex)
}
