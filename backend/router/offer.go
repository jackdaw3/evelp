package router

import (
	"evelp/service"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func offer(c *gin.Context) {
	regionId, err := strconv.Atoi(c.Query("regionId"))
	if err != nil {
		c.AbortWithError(500, err)
	}
	log.Info(c.Query("regionId"))

	scope, err := strconv.ParseFloat(c.Query("scope"), 64)
	if err != nil {
		c.AbortWithError(500, err)
	}
	log.Info(c.Query("scope"))

	corporationId, err := strconv.Atoi(c.Query("corporationId"))
	if err != nil {
		c.AbortWithError(500, err)
	}
	log.Info(c.Query("corporationId"))

	offerService := service.NewOfferSerivce(regionId, float64(scope))

	orders, err := offerService.Offers(corporationId)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.JSON(200, orders)
}
