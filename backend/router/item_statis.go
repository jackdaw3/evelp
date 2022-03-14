package router

import (
	"errors"
	"evelp/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func itemStatis(c *gin.Context) {
	offerId, err := strconv.Atoi(c.Query("offerId"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

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

	tax, err := strconv.ParseFloat(c.Query("tax"), 64)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	materialPrice := c.Query("materialPrice")
	if materialPrice == "" {
		c.AbortWithError(500, errors.New("materialPrice is empty"))
		return
	}

	lang := c.Query("lang")
	if lang == "" {
		c.AbortWithError(500, errors.New("lang is empty"))
		return
	}

	isBuyOrder, err := strconv.ParseBool(c.Query("isBuyOrder"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	itemStatisService := service.NewItemStatisService(offerId, regionId, scope, materialPrice, tax, lang)
	itemStatis, err := itemStatisService.ItemStatis(isBuyOrder)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, itemStatis)
}
