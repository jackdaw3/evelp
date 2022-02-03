package router

import (
	"evelp/model"

	"github.com/gin-gonic/gin"
)

func region(c *gin.Context) {
	factions, err := model.GetRegions()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, factions)
}
