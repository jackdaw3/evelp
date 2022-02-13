package router

import (
	"errors"
	"evelp/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func order(c *gin.Context) {
	regionId, err := strconv.Atoi(c.Query("regionId"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	scope, err := strconv.ParseFloat(c.Query("scope"), 64)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	lang := c.Query("lang")
	if lang == "" {
		c.AbortWithError(500, errors.New("lang is empty"))
		return
	}

	itemId, err := strconv.Atoi(c.Query("itemId"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	isBuyOrder, err := strconv.ParseBool(c.Query("isBuyOrder"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	oderService := service.NewOrderService(itemId, regionId, float64(scope))
	orders, err := oderService.Orders(isBuyOrder, lang)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, orders)
}
