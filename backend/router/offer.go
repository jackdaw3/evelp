package router

import (
	"errors"
	"evelp/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func offer(c *gin.Context) {
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

	days, err := strconv.Atoi(c.Query("days"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	lang := c.Query("lang")
	if lang == "" {
		c.AbortWithError(500, errors.New("lang is empty"))
		return
	}

	corporationId, err := strconv.Atoi(c.Query("corporationId"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	offerService := service.NewOfferSerivce(regionId, float64(scope), days, lang)

	orders, err := offerService.Offers(corporationId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, orders)
}
