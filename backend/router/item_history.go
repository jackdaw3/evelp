package router

import (
	"evelp/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func itemHistory(c *gin.Context) {
	regionId, err := strconv.Atoi(c.Query("regionId"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	itemId, err := strconv.Atoi(c.Query("itemId"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	itemHistoryService := service.NewItemHistoryService(itemId, regionId)
	itemHistorys, err := itemHistoryService.History()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, itemHistorys)
}
