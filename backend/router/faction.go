package router

import (
	"errors"
	"evelp/service"

	"github.com/gin-gonic/gin"
)

func faction(c *gin.Context) {
	lang := c.Query("lang")
	if lang == "" {
		c.AbortWithError(500, errors.New("lang is empty"))
		return
	}

	factionService := service.NewFactionService(lang)
	factions, err := factionService.Factions()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, factions)
}
