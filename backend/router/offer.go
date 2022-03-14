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

	corporationId, err := strconv.Atoi(c.Query("corporationId"))
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

	tax, err := strconv.ParseFloat(c.Query("tax"), 64)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	productPrice := c.Query("productPrice")
	if productPrice == "" {
		c.AbortWithError(500, errors.New("productPrice is empty"))
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

	offerService := service.NewOfferSerivce(regionId, scope, days, productPrice, materialPrice, tax, lang)

	offers, err := offerService.Offers(corporationId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, offers)
}
