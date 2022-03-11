package router

import (
	"errors"
	"evelp/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func corporation(c *gin.Context) {
	corporationId, err := strconv.Atoi(c.Query("corporationId"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	lang := c.Query("lang")
	if lang == "" {
		c.AbortWithError(500, errors.New("lang is empty"))
		return
	}

	corporationService := service.NewCorporationSerivce(corporationId, lang)
	corporation, err := corporationService.Corporation()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, corporation)
}
