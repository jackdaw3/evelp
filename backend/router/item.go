package router

import (
	"errors"
	"evelp/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func item(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Query("itemId"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	lang := c.Query("lang")
	if lang == "" {
		c.AbortWithError(500, errors.New("lang is empty"))
		return
	}

	itemService := service.NewItemService(lang)
	itemDTO, err := itemService.Item(itemId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, itemDTO)
}
